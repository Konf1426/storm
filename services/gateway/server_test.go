package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

type fakeSub struct{}

func (s *fakeSub) Unsubscribe() error { return nil }

type fakeNats struct {
	mu         sync.Mutex
	published  []publishCall
	subCh      chan *nats.Msg
	subscribed chan struct{}
	subjectCh  map[string]chan *nats.Msg
}

type publishCall struct {
	subject string
	data    []byte
}

type noFlushRecorder struct {
	header http.Header
	code   int
	body   bytes.Buffer
}

func (n *noFlushRecorder) Header() http.Header {
	if n.header == nil {
		n.header = http.Header{}
	}
	return n.header
}

func (n *noFlushRecorder) Write(b []byte) (int, error) {
	if n.code == 0 {
		n.code = http.StatusOK
	}
	return n.body.Write(b)
}

func (n *noFlushRecorder) WriteHeader(statusCode int) {
	n.code = statusCode
}

type errNats struct{}

func (errNats) Publish(string, []byte) error { return nil }

func (errNats) ChanSubscribe(string, chan *nats.Msg) (Subscription, error) {
	return nil, errors.New("subscribe failed")
}

func (errNats) IsConnected() bool { return true }

type disconnectedNats struct{}

func (disconnectedNats) Publish(string, []byte) error { return nil }

func (disconnectedNats) ChanSubscribe(string, chan *nats.Msg) (Subscription, error) {
	return &fakeSub{}, nil
}

func (disconnectedNats) IsConnected() bool { return false }

func newFakeNats() *fakeNats {
	return &fakeNats{
		subscribed: make(chan struct{}, 1),
		subjectCh:  make(map[string]chan *nats.Msg),
	}
}

func (f *fakeNats) Publish(subject string, data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.published = append(f.published, publishCall{subject: subject, data: append([]byte(nil), data...)})
	if ch, ok := f.subjectCh[subject]; ok {
		select {
		case ch <- &nats.Msg{Subject: subject, Data: append([]byte(nil), data...)}:
		default:
		}
	}
	return nil
}

func (f *fakeNats) ChanSubscribe(subject string, ch chan *nats.Msg) (Subscription, error) {
	f.subCh = ch
	f.subjectCh[subject] = ch
	select {
	case f.subscribed <- struct{}{}:
	default:
	}
	return &fakeSub{}, nil
}

func (f *fakeNats) IsConnected() bool { return true }

type userRecord struct {
	user         User
	passwordHash string
}

type memStore struct {
	mu          sync.Mutex
	users       map[string]userRecord
	channels    map[int64]Channel
	channelMsgs map[int64][]Message
	refresh     map[string]RefreshToken
	nextChanID  int64
	nextMessage int64
}

func newMemStore() *memStore {
	return &memStore{
		users:       make(map[string]userRecord),
		channels:    make(map[int64]Channel),
		channelMsgs: make(map[int64][]Message),
		refresh:     make(map[string]RefreshToken),
		nextChanID:  1,
		nextMessage: 1,
	}
}

type errorStore struct {
	*memStore
	ensureUserErr error
}

func (e *errorStore) EnsureUser(_ context.Context, _ string) error {
	return e.ensureUserErr
}

func (m *memStore) EnsureUser(_ context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[userID]; ok {
		return nil
	}
	m.users[userID] = userRecord{
		user: User{
			ID:        userID,
			CreatedAt: time.Now(),
		},
	}
	return nil
}

func (m *memStore) CreateUser(_ context.Context, userID, password, displayName string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[userID]; ok {
		return User{}, errors.New("user exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{ID: userID, DisplayName: displayName, CreatedAt: time.Now()}
	m.users[userID] = userRecord{user: user, passwordHash: string(hash)}
	return user, nil
}

func (m *memStore) GetUser(_ context.Context, userID string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.users[userID]
	if !ok {
		return User{}, errors.New("not found")
	}
	return rec.user, nil
}

func (m *memStore) ListUsers(_ context.Context) ([]User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]User, 0, len(m.users))
	for _, rec := range m.users {
		out = append(out, rec.user)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *memStore) UpdateUser(_ context.Context, userID, displayName, password string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.users[userID]
	if !ok {
		return User{}, errors.New("not found")
	}
	if displayName != "" {
		rec.user.DisplayName = displayName
	}
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}
		rec.passwordHash = string(hash)
	}
	m.users[userID] = rec
	return rec.user, nil
}

func (m *memStore) DeleteUser(_ context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, userID)
	return nil
}

func (m *memStore) VerifyUserPassword(_ context.Context, userID, password string) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.users[userID]
	if !ok || rec.passwordHash == "" {
		return User{}, errors.New("invalid user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(rec.passwordHash), []byte(password)); err != nil {
		return User{}, errors.New("invalid password")
	}
	return rec.user, nil
}

func (m *memStore) SaveRefreshToken(_ context.Context, userID, token string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.refresh[token] = RefreshToken{Token: token, UserID: userID, ExpiresAt: expiresAt}
	return nil
}

func (m *memStore) GetRefreshToken(_ context.Context, token string) (RefreshToken, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.refresh[token]
	if !ok {
		return RefreshToken{}, errors.New("not found")
	}
	return rec, nil
}

func (m *memStore) RevokeRefreshToken(_ context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.refresh[token]
	if !ok {
		return nil
	}
	rec.Revoked = true
	m.refresh[token] = rec
	return nil
}

func (m *memStore) CreateChannel(_ context.Context, name, createdBy string) (Channel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[createdBy]; !ok {
		return Channel{}, errors.New("user not found")
	}
	id := m.nextChanID
	m.nextChanID++
	ch := Channel{ID: id, Name: name, CreatedBy: createdBy, CreatedAt: time.Now()}
	m.channels[id] = ch
	return ch, nil
}

func (m *memStore) ListChannels(_ context.Context) ([]Channel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Channel, 0, len(m.channels))
	for _, ch := range m.channels {
		out = append(out, ch)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *memStore) EnsureMember(_ context.Context, channelID int64, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.channels[channelID]; !ok {
		return errors.New("channel not found")
	}
	if _, ok := m.users[userID]; !ok {
		return errors.New("user not found")
	}
	return nil
}

func (m *memStore) SaveChannelMessage(_ context.Context, channelID int64, userID string, payload []byte) (Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg := Message{
		ID:        m.nextMessage,
		ChannelID: channelID,
		UserID:    userID,
		Subject:   channelSubject(channelID),
		Payload:   string(payload),
		CreatedAt: time.Now(),
	}
	m.nextMessage++
	m.channelMsgs[channelID] = append(m.channelMsgs[channelID], msg)
	return msg, nil
}

func (m *memStore) ListMessages(_ context.Context, channelID int64, limit int) ([]Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msgs := m.channelMsgs[channelID]
	if len(msgs) == 0 {
		return []Message{}, nil
	}
	// Return newest first to mirror DB query.
	out := make([]Message, 0, len(msgs))
	for i := len(msgs) - 1; i >= 0; i-- {
		out = append(out, msgs[i])
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (m *memStore) SaveMessage(_ context.Context, _ string, _ []byte) error { return nil }

func (m *memStore) Close() error { return nil }

func TestHealthz(t *testing.T) {
	nc := newFakeNats()
	r := NewRouter(nc, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != "ok" {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
}

func TestPingNatsConnected(t *testing.T) {
	r := NewRouter(newFakeNats(), nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})
	req := httptest.NewRequest(http.MethodGet, "/ping-nats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPingNatsDisconnected(t *testing.T) {
	r := NewRouter(disconnectedNats{}, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})
	req := httptest.NewRequest(http.MethodGet, "/ping-nats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestPublishDefaultSubject(t *testing.T) {
	nc := newFakeNats()
	r := NewRouter(nc, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})

	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBufferString("hello"))
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}

	nc.mu.Lock()
	defer nc.mu.Unlock()
	if len(nc.published) != 1 {
		t.Fatalf("expected 1 publish, got %d", len(nc.published))
	}
	if nc.published[0].subject != defaultSubject {
		t.Fatalf("expected subject %q, got %q", defaultSubject, nc.published[0].subject)
	}
	if string(nc.published[0].data) != "hello" {
		t.Fatalf("unexpected payload: %q", string(nc.published[0].data))
	}
}

func TestPublishEmptyBodyDefaults(t *testing.T) {
	nc := newFakeNats()
	r := NewRouter(nc, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})

	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}

	nc.mu.Lock()
	defer nc.mu.Unlock()
	if len(nc.published) != 1 {
		t.Fatalf("expected 1 publish, got %d", len(nc.published))
	}
	if string(nc.published[0].data) != `{"msg":"hello from gateway"}` {
		t.Fatalf("unexpected payload: %q", string(nc.published[0].data))
	}
}

func TestPublishInvalidSubject(t *testing.T) {
	nc := newFakeNats()
	r := NewRouter(nc, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})

	req := httptest.NewRequest(http.MethodPost, "/publish?subject=bad%20subject", bytes.NewBufferString("hi"))
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPublishTooLarge(t *testing.T) {
	nc := newFakeNats()
	r := NewRouter(nc, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})

	payload := bytes.Repeat([]byte("a"), maxBodyBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Code)
	}
}

type flushRecorder struct {
	header http.Header
	buf    bytes.Buffer
	code   int
}

func newFlushRecorder() *flushRecorder {
	return &flushRecorder{header: make(http.Header)}
}

func (f *flushRecorder) Header() http.Header { return f.header }

func (f *flushRecorder) Write(p []byte) (int, error) {
	if f.code == 0 {
		f.code = http.StatusOK
	}
	return f.buf.Write(p)
}

func (f *flushRecorder) WriteHeader(statusCode int) { f.code = statusCode }

func (f *flushRecorder) Flush() {}

func TestEventsSSE(t *testing.T) {
	nc := newFakeNats()
	r := NewRouter(nc, nil, nil, AuthConfig{Secret: []byte("test-secret"), Enabled: true})

	req := httptest.NewRequest(http.MethodGet, "/events?subject=storm.events", nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	req = req.WithContext(ctx)

	w := newFlushRecorder()

	done := make(chan struct{})
	go func() {
		r.ServeHTTP(w, req)
		close(done)
	}()

	select {
	case <-nc.subscribed:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("subscribe not called")
	}

	nc.subCh <- &nats.Msg{Subject: defaultSubject, Data: []byte("hello")}
	time.Sleep(50 * time.Millisecond)

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("handler did not return")
	}

	body := w.buf.String()
	if !strings.Contains(body, "data: hello") {
		t.Fatalf("expected SSE data, got: %q", body)
	}
}

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, nil, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBufferString("hi"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareAcceptsCookieToken(t *testing.T) {
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, nil, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBufferString("hi"))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: testToken(t, "test-secret")})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
}

func TestReadMessagePayloadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/channels/1/messages", strings.NewReader(`{"payload":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	payload, err := readMessagePayload(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(payload) != "hello" {
		t.Fatalf("expected payload 'hello', got %q", string(payload))
	}
}

func TestReadMessagePayloadEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/channels/1/messages", strings.NewReader(`{"payload":""}`))
	req.Header.Set("Content-Type", "application/json")
	_, err := readMessagePayload(req)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestTokenFromRequestQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws?token=abc", nil)
	if got := tokenFromRequest(req); got != "abc" {
		t.Fatalf("expected token from query")
	}
}

func TestSubjectOrChannel(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws?channel_id=42", nil)
	subject, channelID, err := subjectOrChannel(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if channelID != 42 || subject != "channels.42" {
		t.Fatalf("unexpected subject/channel: %s %d", subject, channelID)
	}

	req = httptest.NewRequest(http.MethodGet, "/ws?subject=storm.events", nil)
	subject, channelID, err = subjectOrChannel(req)
	if err != nil || channelID != 0 || subject != "storm.events" {
		t.Fatalf("unexpected subject/channel")
	}
}

func TestAuthRefreshFlow(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{
		Secret:        []byte("test-secret"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}

	srv := httptest.NewServer(NewRouter(nc, store, nil, auth))
	t.Cleanup(srv.Close)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	_, _ = client.Post(srv.URL+"/auth/register", "application/json", strings.NewReader(`{"user_id":"eve","password":"pass123"}`))
	_, _ = client.Post(srv.URL+"/auth/login", "application/json", strings.NewReader(`{"user_id":"eve","password":"pass123"}`))

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/auth/refresh", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("refresh status %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestAuthLogout(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	token := "refresh-token"
	_ = store.SaveRefreshToken(context.Background(), token, "user-1", time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCookieDomainAndSecureFlags(t *testing.T) {
	cfg := AuthConfig{
		Secret:        []byte("test-secret"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		CookieDomain:  "example.local",
		CookieSecure:  true,
	}

	w := httptest.NewRecorder()
	issueSession(w, cfg, nil, "user-1")
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("expected cookies")
	}
	for _, c := range cookies {
		if c.Domain != "example.local" {
			t.Fatalf("expected domain")
		}
		if !c.Secure {
			t.Fatalf("expected secure cookie")
		}
		if c.MaxAge <= 0 {
			t.Fatalf("expected max age")
		}
	}
}

func TestEnvIntParsing(t *testing.T) {
	t.Setenv("TEST_INT", "12")
	if got := envInt("TEST_INT", 5); got != 12 {
		t.Fatalf("expected 12, got %d", got)
	}
}

func TestEnvBoolParsing(t *testing.T) {
	t.Setenv("TEST_BOOL", "true")
	if !envBool("TEST_BOOL", false) {
		t.Fatalf("expected true")
	}
}

func TestEnvIntFromQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=99", nil)
	if got := envIntFromQuery(req, "limit", 50); got != 99 {
		t.Fatalf("expected 99, got %d", got)
	}
}

func TestReadBodyTooLarge(t *testing.T) {
	tooLarge := bytes.Repeat([]byte("a"), maxBodyBytes+5)
	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBuffer(tooLarge))
	_, err := readBody(nil, req)
	if !errors.Is(err, errPayloadTooLarge) {
		t.Fatalf("expected payload too large")
	}
}

func TestSetCookieExpires(t *testing.T) {
	cfg := AuthConfig{CookieDomain: "example.local"}
	w := httptest.NewRecorder()
	exp := time.Now().Add(1 * time.Hour)
	setCookie(w, "access_token", "abc", exp, cfg)
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("expected cookie")
	}
	if cookies[0].Expires.Before(time.Now()) {
		t.Fatalf("expected future expiry")
	}
}

func TestTokenFromRequestBearer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/publish", nil)
	req.Header.Set("Authorization", "Bearer abc")
	if got := tokenFromRequest(req); got != "abc" {
		t.Fatalf("expected bearer token")
	}
}

func TestSubjectFromRequestDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/publish", nil)
	subject, err := subjectFromRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject != defaultSubject {
		t.Fatalf("expected default subject")
	}
}

func TestParseID(t *testing.T) {
	if _, err := parseID("0"); err == nil {
		t.Fatalf("expected error for id 0")
	}
	if _, err := parseID("abc"); err == nil {
		t.Fatalf("expected error for invalid id")
	}
	if id, err := parseID("12"); err != nil || id != 12 {
		t.Fatalf("expected id 12")
	}
}

func TestTokenFromCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "cookie-token"})
	if got := tokenFromCookie(req, "access_token"); got != "cookie-token" {
		t.Fatalf("expected cookie token")
	}
}

func TestSubjectOrChannelInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws?channel_id=bad", nil)
	_, _, err := subjectOrChannel(req)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
	if !strings.Contains(w.Body.String(), "ok") {
		t.Fatalf("expected body")
	}
}

func TestRefreshTokenInvalid(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{
		Secret:        []byte("test-secret"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}

	srv := httptest.NewServer(NewRouter(nc, store, nil, auth))
	t.Cleanup(srv.Close)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/auth/refresh", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestCorsPreflight(t *testing.T) {
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true, CorsOrigin: "http://example.local"}
	r := NewRouter(nc, newMemStore(), nil, auth)

	req := httptest.NewRequest(http.MethodOptions, "/channels", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://example.local" {
		t.Fatalf("unexpected origin header: %q", got)
	}
}

func TestRegisterInvalidPayload(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegisterMissingFields(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"user_id":"  ","password":" "}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"user_id":"nope","password":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLoginInvalidPayload(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"user_id":`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestChannelsInvalidPayload(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/channels", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestChannelsListSuccess(t *testing.T) {
	store := newMemStore()
	_ = store.EnsureUser(context.Background(), "user-1")
	if _, err := store.CreateChannel(context.Background(), "general", "user-1"); err != nil {
		t.Fatalf("create channel: %v", err)
	}
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(newFakeNats(), store, nil, auth)

	req := httptest.NewRequest(http.MethodGet, "/channels", nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestChannelsListUnauthorized(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodGet, "/channels", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestChannelsEnsureUserError(t *testing.T) {
	store := &errorStore{memStore: newMemStore(), ensureUserErr: errors.New("boom")}
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(newFakeNats(), store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/channels", strings.NewReader(`{"name":"general"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestChannelsStoreNotConfigured(t *testing.T) {
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, nil, nil, auth)

	req := httptest.NewRequest(http.MethodGet, "/channels", nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestChannelsMessagesInvalidID(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/channels/bad/messages", strings.NewReader(`{"payload":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWebSocketChannelFlow(t *testing.T) {
	store := newMemStore()
	_, _ = store.CreateUser(context.Background(), "alice", "pass", "Alice")
	chanRec, _ := store.CreateChannel(context.Background(), "general", "alice")

	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	server := httptest.NewServer(NewRouter(nc, store, nil, auth))
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?channel_id=" + strconv.FormatInt(chanRec.ID, 10)
	header := http.Header{}
	header.Set("Cookie", "access_token="+testToken(t, "test"))
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		t.Fatalf("ws dial: %v", err)
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"user":"alice","message":"hi"}`)); err != nil {
		t.Fatalf("ws write: %v", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ws read: %v", err)
	}
	if !strings.Contains(string(msg), "hi") {
		t.Fatalf("unexpected ws payload: %s", string(msg))
	}
}

func TestWebSocketSubscribeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(wsHandler(errNats{}, nil, nil)))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/?subject=storm.events"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(msg), "subscribe failed") {
		t.Fatalf("unexpected message: %q", string(msg))
	}
}

func TestUsersForbiddenUpdate(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPatch, "/users/other", strings.NewReader(`{"display_name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestUsersGetNotFound(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodGet, "/users/ghost", nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUsersCreateInvalidPayload(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestClearSessionCookies(t *testing.T) {
	cfg := AuthConfig{CookieDomain: "example.local"}
	w := httptest.NewRecorder()
	clearSessionCookies(w, cfg)
	cookies := w.Result().Cookies()
	if len(cookies) < 2 {
		t.Fatalf("expected cookies")
	}
	for _, c := range cookies {
		if c.Expires.After(time.Now()) {
			t.Fatalf("expected expired cookie")
		}
	}
}

func TestPresenceKey(t *testing.T) {
	if got := presenceKey(7); got != "channel:7" {
		t.Fatalf("unexpected presence key: %q", got)
	}
}

func TestAuthFlowAndUsersCRUD(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{
		Secret:        []byte("test-secret"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}

	srv := httptest.NewServer(NewRouter(nc, store, nil, auth))
	t.Cleanup(srv.Close)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	registerBody := `{"user_id":"alice","password":"pass123","display_name":"Alice"}`
	resp, err := client.Post(srv.URL+"/auth/register", "application/json", strings.NewReader(registerBody))
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	resp.Body.Close()

	loginBody := `{"user_id":"alice","password":"pass123"}`
	resp, err = client.Post(srv.URL+"/auth/login", "application/json", strings.NewReader(loginBody))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	resp.Body.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/auth/me", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("me failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("me status %d", resp.StatusCode)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodGet, srv.URL+"/users", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list users status %d", resp.StatusCode)
	}
	resp.Body.Close()

	patchBody := `{"display_name":"Alice Updated"}`
	req, _ = http.NewRequest(http.MethodPatch, srv.URL+"/users/alice", strings.NewReader(patchBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("patch user failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("patch user status %d", resp.StatusCode)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/auth/refresh", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("refresh status %d", resp.StatusCode)
	}
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/users/alice", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("delete user failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete user status %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestChannelsAndMessagesFlow(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{
		Secret:        []byte("test-secret"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}

	srv := httptest.NewServer(NewRouter(nc, store, nil, auth))
	t.Cleanup(srv.Close)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	_, _ = client.Post(srv.URL+"/auth/register", "application/json", strings.NewReader(`{"user_id":"bob","password":"pass123","display_name":"Bob"}`))
	_, _ = client.Post(srv.URL+"/auth/login", "application/json", strings.NewReader(`{"user_id":"bob","password":"pass123"}`))

	createResp, err := client.Post(srv.URL+"/channels", "application/json", strings.NewReader(`{"name":"general"}`))
	if err != nil {
		t.Fatalf("create channel failed: %v", err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create channel status %d", createResp.StatusCode)
	}
	var channel Channel
	if err := json.NewDecoder(createResp.Body).Decode(&channel); err != nil {
		t.Fatalf("decode channel failed: %v", err)
	}
	createResp.Body.Close()

	msgPayload := `{"payload":"{\"user\":\"bob\",\"message\":\"hello\"}"}`
	msgResp, err := client.Post(
		srv.URL+"/channels/"+strconv.FormatInt(channel.ID, 10)+"/messages",
		"application/json",
		strings.NewReader(msgPayload),
	)
	if err != nil {
		t.Fatalf("post message failed: %v", err)
	}
	if msgResp.StatusCode != http.StatusCreated {
		t.Fatalf("post message status %d", msgResp.StatusCode)
	}
	msgResp.Body.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/channels/"+strconv.FormatInt(channel.ID, 10)+"/messages?limit=10", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("list messages failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list messages status %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestClamp(t *testing.T) {
	if got := clamp(5, 1, 10); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
	if got := clamp(-1, 0, 10); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
	if got := clamp(20, 0, 10); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestReadMessagePayloadInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/publish", strings.NewReader(`{"payload":`))
	req.Header.Set("Content-Type", "application/json")
	if _, err := readMessagePayload(req); err == nil {
		t.Fatalf("expected error")
	}
}

func TestReadMessagePayloadEmptyJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/publish", strings.NewReader(`{"payload":"   "}`))
	req.Header.Set("Content-Type", "application/json")
	if _, err := readMessagePayload(req); err == nil {
		t.Fatalf("expected error")
	}
}

func TestStreamSSEUnsupportedWriter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := &noFlushRecorder{}
	streamSSE(rec, req, &fakeNats{}, "storm.events")
	if rec.code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.code)
	}
	if !strings.Contains(rec.body.String(), "streaming unsupported") {
		t.Fatalf("unexpected body: %q", rec.body.String())
	}
}

func TestStreamSSESubscribeError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()
	streamSSE(rec, req, errNats{}, "storm.events")
	if !strings.Contains(rec.Body.String(), "subscribe failed") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestAuthMiddlewareDisabled(t *testing.T) {
	called := false
	handler := authMiddleware(AuthConfig{Enabled: false})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if !called {
		t.Fatalf("expected handler to be called")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestWSHandlerInvalidSubject(t *testing.T) {
	handler := wsHandler(&fakeNats{}, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/ws?subject=bad%20subject", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEnvFallbackAndOverride(t *testing.T) {
	if got := env("STORM_ENV_TEST", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
	t.Setenv("STORM_ENV_TEST", "override")
	if got := env("STORM_ENV_TEST", "fallback"); got != "override" {
		t.Fatalf("expected override, got %q", got)
	}
}

func TestNatsAdapterPublishSubscribe(t *testing.T) {
	opts := &server.Options{
		Host:   "127.0.0.1",
		Port:   -1,
		NoLog:  true,
		NoSigs: true,
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	go ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		t.Fatalf("nats server not ready")
	}
	t.Cleanup(ns.Shutdown)

	// #nosec G402 -- test uses local NATS without TLS.
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	// #nosec G402 -- TLS handled upstream; close is safe.
	defer nc.Close()

	adapter := NewNatsAdapter(nc)
	ch := make(chan *nats.Msg, 1)
	sub, err := adapter.ChanSubscribe("storm.test", ch)
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := adapter.Publish("storm.test", []byte("hello")); err != nil {
		t.Fatalf("publish: %v", err)
	}

	select {
	case msg := <-ch:
		if string(msg.Data) != "hello" {
			t.Fatalf("unexpected payload: %q", string(msg.Data))
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("message not received")
	}

	if !adapter.IsConnected() {
		t.Fatalf("expected adapter connected")
	}
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("unsubscribe: %v", err)
	}
}

func testToken(t *testing.T, secret string) string {
	t.Helper()
	claims := jwt.RegisteredClaims{
		Subject:   "user-1",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("token sign failed: %v", err)
	}
	return signed
}

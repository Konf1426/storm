package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
)

func TestSubjectFromRequestInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/publish?subject=bad%20subject", nil)
	if _, err := subjectFromRequest(req); err == nil {
		t.Fatalf("expected error for invalid subject")
	}
}

func TestSubjectOrChannelDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	subject, channelID, err := subjectOrChannel(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject != defaultSubject || channelID != 0 {
		t.Fatalf("expected default subject, got %q id=%d", subject, channelID)
	}
}

func TestEnvIntInvalidFallback(t *testing.T) {
	t.Setenv("BAD_INT", "not-an-int")
	if got := envInt("BAD_INT", 42); got != 42 {
		t.Fatalf("expected fallback, got %d", got)
	}
}

func TestEnvBoolInvalidFallback(t *testing.T) {
	t.Setenv("BAD_BOOL", "nope")
	if got := envBool("BAD_BOOL", true); got != true {
		t.Fatalf("expected fallback true, got %v", got)
	}
}

func TestCorsDefaultOrigin(t *testing.T) {
	handler := corsMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("unexpected origin header: %q", got)
	}
}

func TestReadBodyWithMaxBytesReader(t *testing.T) {
	payload := bytes.Repeat([]byte("a"), 10)
	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	body, err := readBody(rec, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) != len(payload) {
		t.Fatalf("unexpected body length: %d", len(body))
	}
}

func TestEnvIntFromQueryInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=bad", nil)
	if got := envIntFromQuery(req, "limit", 55); got != 55 {
		t.Fatalf("expected fallback, got %d", got)
	}
}

func TestAuthMeWhenAuthDisabled(t *testing.T) {
	r := NewRouter(newFakeNats(), newMemStore(), nil, AuthConfig{Enabled: false})
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRefreshInvalidSignature(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{
		Secret:        []byte("access"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}

	r := NewRouter(nc, store, nil, auth)

	claims := jwt.RegisteredClaims{
		Subject:   "user-1",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Signed with the wrong secret to trigger signature error.
	signed, err := token.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: signed})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRefreshExpiredToken(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{
		Secret:        []byte("access"),
		RefreshSecret: []byte("refresh-secret"),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}

	r := NewRouter(nc, store, nil, auth)

	claims := jwt.RegisteredClaims{
		Subject:   "user-1",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(auth.RefreshSecret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	_ = store.SaveRefreshToken(context.Background(), "user-1", signed, time.Now().Add(-time.Minute))

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: signed})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestWriteSSE(t *testing.T) {
	var buf bytes.Buffer
	if err := writeSSE(&buf, "hello"); err != nil {
		t.Fatalf("writeSSE error: %v", err)
	}
	if got := buf.String(); got != "data: hello\n\n" {
		t.Fatalf("unexpected SSE payload: %q", got)
	}
}

func TestWriteSSEError(t *testing.T) {
	if err := writeSSE(errWriter{}, "hello"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSignToken(t *testing.T) {
	secret := []byte("secret")
	tokenStr, exp := signToken(secret, "user-1", time.Minute)
	if exp.Before(time.Now()) {
		t.Fatalf("expected future expiration")
	}
	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("token parse failed: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Fatalf("unexpected subject: %q", claims.Subject)
	}
}

func TestWriteJSONAdditional(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusCreated, map[string]string{"status": "ok"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
	if body := rec.Body.String(); body == "" || body[0] != '{' {
		t.Fatalf("unexpected body: %q", body)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content-type: %q", got)
	}
}

func TestTokenFromCookieMissing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	if got := tokenFromCookie(req, "missing"); got != "" {
		t.Fatalf("expected empty token, got %q", got)
	}
}

func TestTokenFromRequestUsesBearerFirst(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "cookie-token"})
	if got := tokenFromRequest(req); got != "header-token" {
		t.Fatalf("expected header token, got %q", got)
	}
}

func TestTokenFromRequestCookieFallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "cookie-token"})
	if got := tokenFromRequest(req); got != "cookie-token" {
		t.Fatalf("expected cookie token, got %q", got)
	}
}

func TestTokenFromRequestQueryFallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws?token=query-token", nil)
	if got := tokenFromRequest(req); got != "query-token" {
		t.Fatalf("expected query token, got %q", got)
	}
}

func TestSubjectFromRequestDefaultAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/publish", nil)
	got, err := subjectFromRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != defaultSubject {
		t.Fatalf("expected default subject, got %q", got)
	}
}

type stubSub struct{}

func (stubSub) Unsubscribe() error { return nil }

type stubNats struct {
	SubscribeErr error
	Msgs         chan *nats.Msg
}

func (s stubNats) Publish(string, []byte) error { return nil }

func (s stubNats) ChanSubscribe(string, chan *nats.Msg) (Subscription, error) {
	if s.SubscribeErr != nil {
		return nil, s.SubscribeErr
	}
	return stubSub{}, nil
}

func (s stubNats) IsConnected() bool { return true }

type sendingNats struct {
	Msg     *nats.Msg
	SendNil bool
}

func (s sendingNats) Publish(string, []byte) error { return nil }

func (s sendingNats) ChanSubscribe(_ string, ch chan *nats.Msg) (Subscription, error) {
	msg := s.Msg
	sendNil := s.SendNil
	go func() {
		if sendNil {
			ch <- nil
			return
		}
		if msg != nil {
			ch <- msg
		}
	}()
	return stubSub{}, nil
}

func (s sendingNats) IsConnected() bool { return true }

type flushRecorderAdditional struct {
	code   int
	header http.Header
	body   bytes.Buffer
}

func newFlushRecorderAdditional() *flushRecorderAdditional {
	return &flushRecorderAdditional{header: make(http.Header)}
}

func (f *flushRecorderAdditional) Header() http.Header { return f.header }

func (f *flushRecorderAdditional) Write(p []byte) (int, error) {
	if f.code == 0 {
		f.code = http.StatusOK
	}
	return f.body.Write(p)
}

func (f *flushRecorderAdditional) WriteHeader(code int) { f.code = code }

func (f *flushRecorderAdditional) Flush() {}

type noFlushRecorderAdditional struct {
	code   int
	header http.Header
	body   bytes.Buffer
}

func newNoFlushRecorderAdditional() *noFlushRecorderAdditional {
	return &noFlushRecorderAdditional{header: make(http.Header)}
}

func (n *noFlushRecorderAdditional) Header() http.Header { return n.header }

func (n *noFlushRecorderAdditional) Write(p []byte) (int, error) {
	if n.code == 0 {
		n.code = http.StatusOK
	}
	return n.body.Write(p)
}

func (n *noFlushRecorderAdditional) WriteHeader(code int) { n.code = code }

func TestStreamSSEUnsupportedWriterAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := newNoFlushRecorderAdditional()
	streamSSE(rec, req, stubNats{}, "storm.events")
	if rec.code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.code)
	}
	if !strings.Contains(rec.body.String(), "streaming unsupported") {
		t.Fatalf("unexpected body: %q", rec.body.String())
	}
}

func TestStreamSSESubscribeErrorAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := newFlushRecorderAdditional()
	streamSSE(rec, req, stubNats{SubscribeErr: errors.New("boom")}, "storm.events")
	if rec.code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.code)
	}
	if !strings.Contains(rec.body.String(), "subscribe failed") {
		t.Fatalf("unexpected body: %q", rec.body.String())
	}
}

func TestStreamSSEHeartbeatAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	rec := newFlushRecorderAdditional()
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	streamSSE(rec, req, stubNats{}, "storm.events")
	if !strings.Contains(rec.body.String(), ": stream ready") {
		t.Fatalf("expected stream ready prelude")
	}
}

func TestWriteSSEMultiline(t *testing.T) {
	var buf bytes.Buffer
	if err := writeSSE(&buf, "a\nb"); err != nil {
		t.Fatalf("writeSSE error: %v", err)
	}
	if got := buf.String(); got != "data: a\ndata: b\n\n" {
		t.Fatalf("unexpected SSE payload: %q", got)
	}
}

func TestStreamSSEMessageAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	rec := newFlushRecorderAdditional()
	done := make(chan struct{})
	go func() {
		streamSSE(rec, req, sendingNats{Msg: &nats.Msg{Data: []byte("hello")}}, "storm.events")
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("stream did not stop")
	}
	if !strings.Contains(rec.body.String(), "data: hello") {
		t.Fatalf("expected SSE message, got %q", rec.body.String())
	}
}

func TestStreamSSENilMessageAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	rec := newFlushRecorderAdditional()
	done := make(chan struct{})
	go func() {
		streamSSE(rec, req, sendingNats{SendNil: true}, "storm.events")
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("stream did not stop")
	}
	if !strings.Contains(rec.body.String(), ": stream ready") {
		t.Fatalf("expected stream ready prelude")
	}
}

type errReadCloser struct{}

func (errReadCloser) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (errReadCloser) Close() error             { return nil }

func TestReadBodyNoWriterAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/publish", io.NopCloser(strings.NewReader("hello")))
	body, err := readBody(nil, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("unexpected body: %q", string(body))
	}
}

func TestReadBodyErrorAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/publish", errReadCloser{})
	if _, err := readBody(nil, req); err == nil {
		t.Fatalf("expected error")
	}
}

func TestUserFromContextInvalidTypeAdditional(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxUserIDKey{}, 123)
	if got := userFromContext(ctx); got != "" {
		t.Fatalf("expected empty user, got %q", got)
	}
}

func TestTokenFromRequestBearerEmptyAdditional(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws?token=query", nil)
	req.Header.Set("Authorization", "Bearer ")
	if got := tokenFromRequest(req); got != "" {
		t.Fatalf("expected empty token, got %q", got)
	}
}

type errStore struct{}

func (errStore) EnsureUser(context.Context, string) error { return nil }
func (errStore) CreateUser(context.Context, string, string, string) (User, error) {
	return User{}, nil
}
func (errStore) GetUser(context.Context, string) (User, error) { return User{}, nil }
func (errStore) ListUsers(context.Context) ([]User, error)     { return nil, nil }
func (errStore) UpdateUser(context.Context, string, string, string) (User, error) {
	return User{}, nil
}
func (errStore) DeleteUser(context.Context, string) error { return nil }
func (errStore) VerifyUserPassword(context.Context, string, string) (User, error) {
	return User{}, nil
}
func (errStore) SaveRefreshToken(context.Context, string, string, time.Time) error {
	return errors.New("save refresh failed")
}
func (errStore) GetRefreshToken(context.Context, string) (RefreshToken, error) { return RefreshToken{}, nil }
func (errStore) RevokeRefreshToken(context.Context, string) error               { return nil }
func (errStore) CreateChannel(context.Context, string, string) (Channel, error) {
	return Channel{}, nil
}
func (errStore) ListChannels(context.Context) ([]Channel, error) { return nil, nil }
func (errStore) EnsureMember(context.Context, int64, string) error {
	return nil
}
func (errStore) SaveChannelMessage(context.Context, int64, string, []byte) (Message, error) {
	return Message{}, nil
}
func (errStore) ListMessages(context.Context, int64, int) ([]Message, error) { return nil, nil }
func (errStore) SaveMessage(context.Context, string, []byte) error           { return nil }
func (errStore) Close() error                                                { return nil }

func TestIssueSessionStoreErrorAdditional(t *testing.T) {
	cfg := AuthConfig{
		Secret:        []byte("test-secret"),
		RefreshSecret: []byte("refresh-secret"),
		AccessTTL:     time.Minute,
		RefreshTTL:    2 * time.Minute,
	}
	w := httptest.NewRecorder()
	issueSession(w, cfg, errStore{}, "user-1")
	cookies := w.Result().Cookies()
	if len(cookies) < 2 {
		t.Fatalf("expected cookies to be set")
	}
}

func TestReadBodyTooLargeAdditional(t *testing.T) {
	payload := bytes.Repeat([]byte("a"), maxBodyBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBuffer(payload))
	if _, err := readBody(nil, req); err == nil {
		t.Fatalf("expected payload too large error")
	}
}

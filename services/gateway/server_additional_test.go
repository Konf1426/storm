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

func TestEnvIntValid(t *testing.T) {
	t.Setenv("GOOD_INT", "42")
	if got := envInt("GOOD_INT", 0); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestEnvBoolInvalidFallback(t *testing.T) {
	t.Setenv("BAD_BOOL", "nope")
	if got := envBool("BAD_BOOL", true); got != true {
		t.Fatalf("expected fallback true, got %v", got)
	}
}

func TestEnvBoolValid(t *testing.T) {
	t.Setenv("GOOD_BOOL", "true")
	if got := envBool("GOOD_BOOL", false); got != true {
		t.Fatalf("expected true, got %v", got)
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

func TestWorkerPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	store := errStore{}
	StartWorkerPool(ctx, store, 1)

	asyncTaskQueue <- asyncTask{
		taskType:  taskSaveMessage,
		channelID: 1,
		userID:    "user",
		payload:   []byte("test"),
	}

	asyncTaskQueue <- asyncTask{
		taskType:  taskSaveRefreshToken,
		userID:    "user",
		token:     "token",
		expiresAt: time.Now(),
	}

	time.Sleep(50 * time.Millisecond) // Let workers process tasks
	cancel()
	time.Sleep(50 * time.Millisecond) // Let workers exit cleanly
}

func TestRateLimiterGet(t *testing.T) {
	rl := newRateLimiter(10, 10)
	l := rl.get("127.0.0.1")
	if l == nil {
		t.Fatalf("expected non-nil limiter")
	}

	l2 := rl.get("127.0.0.1")
	if l != l2 {
		t.Fatalf("expected exactly the same limiter on second get")
	}

	// simulate old entry for cleanup coverage
	rl.mu.Lock()
	e := rl.limiters["127.0.0.1"]
	e.lastSeen = time.Now().Add(-11 * time.Minute)
	rl.mu.Unlock()
}

func TestAuthRateLimitMiddleware(t *testing.T) {
	t.Setenv("AUTH_RATE_LIMIT_ENABLED", "true")
	rl := newRateLimiter(0, 0) // Deny all requests
	handler := rl.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Attempt request while enabled and rate limited
	req := httptest.NewRequest("GET", "/auth", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", rec.Code)
	}

	// Disable rate limiting with ENV var
	t.Setenv("AUTH_RATE_LIMIT_ENABLED", "false")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK since disabled, got %d", rec2.Code)
	}
}

func TestSanitize(t *testing.T) {
	input := "user\nname\r"
	expected := "username"
	if got := sanitize(input); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

type mockNats struct {
	connected bool
}

func (m *mockNats) Publish(subject string, data []byte) error { return nil }
func (m *mockNats) ChanSubscribe(subject string, ch chan *nats.Msg) (Subscription, error) {
	return nil, nil
}
func (m *mockNats) IsConnected() bool { return m.connected }

func TestHealthzAndPingNats(t *testing.T) {
	r := NewRouter(&mockNats{connected: false}, nil, nil, AuthConfig{})

	reqHealth := httptest.NewRequest("GET", "/healthz", nil)
	recHealth := httptest.NewRecorder()
	r.ServeHTTP(recHealth, reqHealth)
	if recHealth.Code != http.StatusOK {
		t.Fatalf("healthz: expected 200, got %d", recHealth.Code)
	}

	reqPing := httptest.NewRequest("GET", "/ping-nats", nil)
	recPing := httptest.NewRecorder()
	r.ServeHTTP(recPing, reqPing)
	if recPing.Code != http.StatusServiceUnavailable {
		t.Fatalf("ping-nats: expected 503, got %d", recPing.Code)
	}

	r2 := NewRouter(&mockNats{connected: true}, nil, nil, AuthConfig{})
	recPing2 := httptest.NewRecorder()
	r2.ServeHTTP(recPing2, reqPing)
	if recPing2.Code != http.StatusOK {
		t.Fatalf("ping-nats: expected 200, got %d", recPing2.Code)
	}
}

func TestRefreshMissingToken(t *testing.T) {
	r := NewRouter(&mockNats{}, nil, nil, AuthConfig{})
	req := httptest.NewRequest("POST", "/auth/refresh", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable { // store not configured
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	r2 := NewRouter(&mockNats{}, errStore{}, nil, AuthConfig{})
	req2 := httptest.NewRequest("POST", "/auth/refresh", nil)
	rec2 := httptest.NewRecorder()
	r2.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusUnauthorized { // missing refresh token
		t.Fatalf("expected 401 for missing token, got %d", rec2.Code)
	}
}

func TestLogout(t *testing.T) {
	r := NewRouter(&mockNats{}, errStore{}, nil, AuthConfig{})
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

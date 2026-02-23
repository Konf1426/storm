package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// ---------------------------------------------------------------------------
// Rate Limiter tests
// ---------------------------------------------------------------------------

func TestRateLimiterAllowsWithinBurst(t *testing.T) {
	rl := newRateLimiter(rate.Every(time.Hour), 2) // very slow refill, burst of 2
	handler := rl.middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, w.Code)
		}
	}
}

func TestRateLimiterRejects429(t *testing.T) {
	rl := newRateLimiter(rate.Every(time.Hour), 1) // burst of 1
	handler := rl.middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request should succeed.
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", w.Code)
	}

	// Second request from same IP should be rate limited.
	req = httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: expected 429, got %d", w.Code)
	}
}

func TestRateLimiterDifferentIPsAreIndependent(t *testing.T) {
	rl := newRateLimiter(rate.Every(time.Hour), 1) // burst of 1
	handler := rl.middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Exhaust limit for IP-A
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// IP-B should still be allowed.
	req = httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.2:5678"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("different IP should succeed, got %d", w.Code)
	}
}

func TestRateLimiterRemoteAddrWithoutPort(t *testing.T) {
	rl := newRateLimiter(rate.Every(time.Hour), 1)
	handler := rl.middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// RemoteAddr without port (edge case: SplitHostPort will fail, fallback to raw string)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Second request from same "ip" should be rate limited.
	req = httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Auth Middleware – expired access token
// ---------------------------------------------------------------------------

func TestAuthMiddlewareRejectsExpiredToken(t *testing.T) {
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(newFakeNats(), nil, nil, auth)

	claims := jwt.RegisteredClaims{
		Subject:   "user-1",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // already expired
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(auth.Secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBufferString("hi"))
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Auth Middleware – wrong signing algorithm (none / RS256)
// ---------------------------------------------------------------------------

func TestAuthMiddlewareRejectsBadAlgo(t *testing.T) {
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(newFakeNats(), nil, nil, auth)

	// Sign with wrong secret (simulates invalid signature)
	claims := jwt.RegisteredClaims{
		Subject:   "user-1",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("wrong-secret"))

	req := httptest.NewRequest(http.MethodPost, "/publish", bytes.NewBufferString("hi"))
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for bad signature, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Auth – register duplicate user
// ---------------------------------------------------------------------------

func TestRegisterDuplicateUser(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true, AccessTTL: time.Minute, RefreshTTL: time.Hour}
	r := NewRouter(nc, store, nil, auth)

	body := `{"user_id":"alice","password":"pass123"}`

	// First registration should succeed.
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.RemoteAddr = "10.10.10.1:1111"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("first register: expected 2xx, got %d", w.Code)
	}

	// Second registration with same user_id should fail.
	req = httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.RemoteAddr = "10.10.10.2:2222" // different IP to avoid rate limiter
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusCreated || w.Code == http.StatusOK {
		t.Fatalf("duplicate register should fail, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Auth – register without store returns 503
// ---------------------------------------------------------------------------

func TestRegisterNoStore(t *testing.T) {
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, nil, nil, auth) // nil store

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"user_id":"bob","password":"pass123"}`))
	req.RemoteAddr = "10.20.20.1:3333"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Auth – login without store returns 503
// ---------------------------------------------------------------------------

func TestLoginNoStore(t *testing.T) {
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, nil, nil, auth) // nil store

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"user_id":"bob","password":"pass123"}`))
	req.RemoteAddr = "10.30.30.1:4444"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// Auth – /auth/me returns user info with valid token
// ---------------------------------------------------------------------------

func TestAuthMeReturnsUser(t *testing.T) {
	store := newMemStore()
	nc := newFakeNats()
	auth := AuthConfig{Secret: []byte("test-secret"), RefreshSecret: []byte("refresh"), Enabled: true}
	r := NewRouter(nc, store, nil, auth)

	// Create a user so /auth/me can look them up.
	_ = store.EnsureUser(nil, "user-1")

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.RemoteAddr = "10.40.40.1:5555"
	req.Header.Set("Authorization", "Bearer "+testToken(t, "test-secret"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "user-1") {
		t.Fatalf("expected body to contain user-1, got %q", w.Body.String())
	}
}

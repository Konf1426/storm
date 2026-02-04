package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

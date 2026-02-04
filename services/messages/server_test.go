package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestHealthz(t *testing.T) {
	r := NewRouter()
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

func TestMetricsEndpoint(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "go_gc_duration_seconds") {
		t.Fatalf("expected prometheus metrics output")
	}
}

func TestEnvFallback(t *testing.T) {
	if got := env("MESSAGES_ADDR", ":9999"); got != ":9999" {
		t.Fatalf("expected fallback, got %q", got)
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("MESSAGES_ADDR", ":1234")
	if got := env("MESSAGES_ADDR", ":9999"); got != ":1234" {
		t.Fatalf("expected override, got %q", got)
	}
}

func TestHealthzIsFast(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	start := time.Now()
	r.ServeHTTP(w, req)
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("healthz too slow")
	}
}

type fakeNatsConn struct {
	subscribeErr error
	flushErr     error
	lastErr      error
}

func (f *fakeNatsConn) Subscribe(_ string, _ nats.MsgHandler) (*nats.Subscription, error) {
	if f.subscribeErr != nil {
		return nil, f.subscribeErr
	}
	return &nats.Subscription{}, nil
}

func (f *fakeNatsConn) Flush() error { return f.flushErr }

func (f *fakeNatsConn) LastError() error { return f.lastErr }

func (f *fakeNatsConn) Close() {}

func TestRunMessagesSubscribeError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := runMessages(ctx, &fakeNatsConn{subscribeErr: errors.New("subscribe failed")}, "s", ":0")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunMessagesFlushError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := runMessages(ctx, &fakeNatsConn{flushErr: errors.New("flush failed")}, "s", ":0")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunMessagesLastError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := runMessages(ctx, &fakeNatsConn{lastErr: errors.New("last error")}, "s", ":0")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunMessagesCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	err := runMessages(ctx, &fakeNatsConn{}, "s", ":0")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRunMessagesWithConnectorError(t *testing.T) {
	ctx := context.Background()
	err := runMessagesWithConnector(ctx, func(string) (NatsConn, error) {
		return nil, errors.New("dial error")
	}, "nats://x", "s", ":0")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunMainUsesConnector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	called := false
	err := runMain(ctx, func(string) (NatsConn, error) {
		called = true
		return &fakeNatsConn{}, nil
	}, "nats://x", "s", ":0")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !called {
		t.Fatalf("expected connector to be called")
	}
}

func TestRunEntryUsesDeps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	notifyCalled := false
	connectCalled := false

	t.Setenv("MESSAGES_ADDR", ":0")

	deps := runtimeDeps{
		NotifyContext: func(context.Context, ...os.Signal) (context.Context, context.CancelFunc) {
			notifyCalled = true
			return ctx, cancel
		},
		Connect: func(string) (NatsConn, error) {
			connectCalled = true
			return &fakeNatsConn{}, nil
		},
		Logf: func(string, ...any) {},
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	if err := runEntry(deps); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !notifyCalled || !connectCalled {
		t.Fatalf("expected deps to be called")
	}
}

func TestRunEntryConnectError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Setenv("MESSAGES_ADDR", ":0")
	t.Setenv("PPROF_ADDR", "")

	deps := runtimeDeps{
		NotifyContext: func(context.Context, ...os.Signal) (context.Context, context.CancelFunc) {
			return ctx, cancel
		},
		Connect: func(string) (NatsConn, error) {
			return nil, errors.New("dial error")
		},
		// Logf left nil to cover default.
	}

	if err := runEntry(deps); err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunEntryWithPprofEnabled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Setenv("MESSAGES_ADDR", ":0")
	t.Setenv("PPROF_ADDR", "127.0.0.1:0")

	deps := runtimeDeps{
		NotifyContext: func(context.Context, ...os.Signal) (context.Context, context.CancelFunc) {
			return ctx, cancel
		},
		Connect: func(string) (NatsConn, error) {
			return &fakeNatsConn{}, nil
		},
		Logf: func(string, ...any) {},
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	if err := runEntry(deps); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

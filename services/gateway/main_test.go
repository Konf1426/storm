package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type dummyStore struct{}

func (d dummyStore) EnsureUser(context.Context, string) error { return nil }
func (d dummyStore) CreateUser(context.Context, string, string, string) (User, error) {
	return User{}, nil
}
func (d dummyStore) GetUser(context.Context, string) (User, error) { return User{}, nil }
func (d dummyStore) ListUsers(context.Context) ([]User, error)     { return nil, nil }
func (d dummyStore) UpdateUser(context.Context, string, string, string) (User, error) {
	return User{}, nil
}
func (d dummyStore) DeleteUser(context.Context, string) error { return nil }
func (d dummyStore) VerifyUserPassword(context.Context, string, string) (User, error) {
	return User{}, nil
}
func (d dummyStore) SaveRefreshToken(context.Context, string, string, time.Time) error { return nil }
func (d dummyStore) GetRefreshToken(context.Context, string) (RefreshToken, error) {
	return RefreshToken{}, nil
}
func (d dummyStore) RevokeRefreshToken(context.Context, string) error { return nil }
func (d dummyStore) CreateChannel(context.Context, string, string) (Channel, error) {
	return Channel{}, nil
}
func (d dummyStore) ListChannels(context.Context) ([]Channel, error)   { return nil, nil }
func (d dummyStore) EnsureMember(context.Context, int64, string) error { return nil }
func (d dummyStore) SaveChannelMessage(context.Context, int64, string, []byte) (Message, error) {
	return Message{}, nil
}
func (d dummyStore) ListMessages(context.Context, int64, int) ([]Message, error) { return nil, nil }
func (d dummyStore) SaveMessage(context.Context, string, []byte) error           { return nil }
func (d dummyStore) Close() error                                                { return nil }

type dummyPresence struct{}

func (d dummyPresence) Incr(context.Context, string) error { return nil }
func (d dummyPresence) Decr(context.Context, string) error { return nil }
func (d dummyPresence) Close() error                       { return nil }

func TestConnectPostgresFailure(t *testing.T) {
	ctx := context.Background()
	_, err := connectPostgres(ctx, "postgres://bad:bad@127.0.0.1:1/bad?sslmode=disable", 1, 1*time.Millisecond)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestConnectRedisFailure(t *testing.T) {
	ctx := context.Background()
	_, err := connectRedis(ctx, "127.0.0.1:1", "", 0, 1, 1*time.Millisecond)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestConnectRedisSuccess(t *testing.T) {
	ctx := context.Background()
	srv := miniredis.RunT(t)
	_, err := connectRedis(ctx, srv.Addr(), "", 0, 1, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestRunMainUsesDeps(t *testing.T) {
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

	t.Setenv("NATS_URL", ns.ClientURL())
	t.Setenv("GATEWAY_ADDR", ":0")

	var (
		pgCalled     bool
		redisCalled  bool
		listenCalled bool
	)
	listenErr := errors.New("listen error")
	deps := runtimeDeps{
		// #nosec G402 -- test uses local NATS without TLS.
		NatsConnect: nats.Connect,
		ConnectPostgres: func(context.Context, string, int, time.Duration) (Store, error) {
			pgCalled = true
			return dummyStore{}, nil
		},
		ConnectRedis: func(context.Context, string, string, int, int, time.Duration) (Presence, error) {
			redisCalled = true
			return dummyPresence{}, nil
		},
		ListenAndServe: func(string, http.Handler) error {
			listenCalled = true
			return listenErr
		},
	}

	if err := runMain(deps); !errors.Is(err, listenErr) {
		t.Fatalf("expected listen error, got %v", err)
	}
	if !pgCalled || !redisCalled || !listenCalled {
		t.Fatalf("expected deps to be called")
	}
}

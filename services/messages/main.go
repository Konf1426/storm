package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if err := runEntry(runtimeDeps{}); err != nil {
		log.Fatalf("messages failed: %v", err)
	}
}

type runtimeDeps struct {
	NotifyContext func(context.Context, ...os.Signal) (context.Context, context.CancelFunc)
	Connect       func(string) (NatsConn, error)
	Logf          func(string, ...any)
}

func runEntry(deps runtimeDeps) error {
	if deps.NotifyContext == nil {
		deps.NotifyContext = signal.NotifyContext
	}
	if deps.Connect == nil {
		deps.Connect = func(url string) (NatsConn, error) {
			return nats.Connect(
				url,
				nats.Name("storm-messages"),
				nats.Timeout(3*time.Second),
				nats.PingInterval(20*time.Second),
				nats.MaxPingsOutstanding(3),
			)
		}
	}
	if deps.Logf == nil {
		deps.Logf = log.Printf
	}

	addr := env("MESSAGES_ADDR", ":8081")
	natsURL := env("NATS_URL", "nats://localhost:4222")
	subject := env("SUBJECT", "storm.events")

	ctx, stop := deps.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	deps.Logf("messages listening on %s (nats=%s subject=%s)", addr, natsURL, subject)
	return runMain(ctx, deps.Connect, natsURL, subject, addr)
}

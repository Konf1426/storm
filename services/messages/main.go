package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// NATS
	natsURL := env("NATS_URL", "nats://localhost:4222")
	subject := env("SUBJECT", "storm.events")

	nc, err := nats.Connect(
		natsURL,
		nats.Name("storm-messages"),
		nats.Timeout(3*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(3),
	)
	if err != nil {
		log.Fatalf("nats connect failed: %v", err)
	}
	defer nc.Close()

	// Subscribe
	_, err = nc.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("received subject=%s bytes=%d payload=%s", m.Subject, len(m.Data), string(m.Data))
	})
	if err != nil {
		log.Fatalf("subscribe failed: %v", err)
	}
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatalf("nats flush error: %v", err)
	}

	addr := env("MESSAGES_ADDR", ":8081")
	srv := &http.Server{Addr: addr, Handler: NewRouter()}

	go func() {
		log.Printf("messages listening on %s (nats=%s subject=%s)", addr, natsURL, subject)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")
	_ = srv.Close()
}

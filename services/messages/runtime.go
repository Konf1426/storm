package main

import (
	"context"
	"net/http"

	"github.com/nats-io/nats.go"
)

type NatsConn interface {
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	Flush() error
	LastError() error
	Close()
}

func runMessages(ctx context.Context, nc NatsConn, subject, addr string) error {
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {})
	if err != nil {
		return err
	}

	if err := nc.Flush(); err != nil {
		return err
	}
	if err := nc.LastError(); err != nil {
		return err
	}

	srv := &http.Server{Addr: addr, Handler: NewRouter()}

	errCh := make(chan error, 1)
	go func() {
		//gosec:ignore G402 -- TLS handled at ingress; service is internal.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		_ = srv.Close()
		return nil
	case err := <-errCh:
		return err
	}
}

func runMessagesWithConnector(ctx context.Context, connect func(string) (NatsConn, error), natsURL, subject, addr string) error {
	nc, err := connect(natsURL)
	if err != nil {
		return err
	}
	defer nc.Close()
	return runMessages(ctx, nc, subject, addr)
}

func runMain(ctx context.Context, connect func(string) (NatsConn, error), natsURL, subject, addr string) error {
	return runMessagesWithConnector(ctx, connect, natsURL, subject, addr)
}

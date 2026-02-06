package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if err := runMain(runtimeDeps{}); err != nil {
		log.Fatal(err)
	}
}

type runtimeDeps struct {
	NatsConnect     func(string, ...nats.Option) (*nats.Conn, error)
	ConnectPostgres func(context.Context, string, int, time.Duration) (Store, error)
	ConnectRedis    func(context.Context, string, string, int, int, time.Duration) (Presence, error)
	ListenAndServe  func(string, http.Handler) error
}

func runMain(deps runtimeDeps) error {
	if deps.NatsConnect == nil {
		deps.NatsConnect = nats.Connect
	}
	if deps.ConnectPostgres == nil {
		deps.ConnectPostgres = connectPostgres
	}
	if deps.ConnectRedis == nil {
		deps.ConnectRedis = connectRedis
	}
	if deps.ListenAndServe == nil {
		deps.ListenAndServe = http.ListenAndServe
	}

	natsURL := env("NATS_URL", "nats://localhost:4222")
	//gosec:ignore G402 -- TLS is terminated at ingress in prod; dev uses plaintext.
	nc, err := deps.NatsConnect(
		natsURL,
		nats.Name("storm-gateway"),
		nats.Timeout(3*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(3),
	)
	if err != nil {
		return err
	}
	//gosec:ignore G402 -- TLS handled upstream; close is safe.
	defer nc.Close()

	ctx := context.Background()
	pgDSN := env("POSTGRES_DSN", "postgres://storm:storm@postgres:5432/storm?sslmode=disable")
	store, err := deps.ConnectPostgres(ctx, pgDSN, 10, 2*time.Second)
	if err != nil {
		return err
	}
	defer func() {
		_ = store.Close()
	}()

	redisAddr := env("REDIS_ADDR", "redis:6379")
	redisPassword := env("REDIS_PASSWORD", "")
	redisDB := envInt("REDIS_DB", 0)
	presence, err := deps.ConnectRedis(ctx, redisAddr, redisPassword, redisDB, 10, 2*time.Second)
	if err != nil {
		return err
	}
	defer func() {
		_ = presence.Close()
	}()

	jwtSecret := env("JWT_SECRET", "dev-secret")
	jwtRefreshSecret := env("JWT_REFRESH_SECRET", jwtSecret)
	auth := AuthConfig{
		Secret:        []byte(jwtSecret),
		RefreshSecret: []byte(jwtRefreshSecret),
		Enabled:       true,
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		CookieDomain:  env("COOKIE_DOMAIN", ""),
		CookieSecure:  envBool("COOKIE_SECURE", false),
		CorsOrigin:    env("CORS_ORIGIN", "http://localhost:5173"),
	}

	addr := env("GATEWAY_ADDR", ":8080")
	log.Printf("gateway listening on %s (nats=%s)", addr, natsURL)
	if pprofAddr := env("PPROF_ADDR", ""); pprofAddr != "" {
		go func() {
			log.Printf("pprof listening on %s", pprofAddr)
			//gosec:ignore G402 -- dev-only pprof, TLS upstream.
			if err := http.ListenAndServe(pprofAddr, nil); err != nil {
				log.Printf("pprof error: %v", err)
			}
		}()
	}
	//gosec:ignore G402 -- TLS at ingress.
	return deps.ListenAndServe(addr, NewRouter(NewNatsAdapter(nc), store, presence, auth))
}

func connectPostgres(ctx context.Context, dsn string, attempts int, delay time.Duration) (Store, error) {
	var lastErr error
	for i := 1; i <= attempts; i++ {
		store, err := NewPostgresStore(ctx, dsn)
		if err == nil {
			return store, nil
		}
		lastErr = err
		log.Printf("postgres connect attempt %d/%d failed: %v", i, attempts, err)
		time.Sleep(delay)
	}
	return nil, lastErr
}

func connectRedis(ctx context.Context, addr, password string, db int, attempts int, delay time.Duration) (Presence, error) {
	var lastErr error
	for i := 1; i <= attempts; i++ {
		presence, err := NewRedisPresence(ctx, addr, password, db)
		if err == nil {
			return presence, nil
		}
		lastErr = err
		log.Printf("redis connect attempt %d/%d failed: %v", i, attempts, err)
		time.Sleep(delay)
	}
	return nil, lastErr
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultSubject = "storm.events"
	maxBodyBytes   = 1 << 20
)

var subjectRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
var errPayloadTooLarge = errors.New("payload too large")

// NatsClient is the minimal interface needed by the HTTP handlers.
type NatsClient interface {
	Publish(subject string, data []byte) error
	ChanSubscribe(subject string, ch chan *nats.Msg) (Subscription, error)
	IsConnected() bool
}

// Subscription is a minimal wrapper for NATS subscriptions.
type Subscription interface {
	Unsubscribe() error
}

// Store persists users/channels/messages and tokens.
type Store interface {
	EnsureUser(ctx context.Context, userID string) error
	CreateUser(ctx context.Context, userID, password, displayName string) (User, error)
	GetUser(ctx context.Context, userID string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, userID, displayName, password string) (User, error)
	DeleteUser(ctx context.Context, userID string) error
	VerifyUserPassword(ctx context.Context, userID, password string) (User, error)
	SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	CreateChannel(ctx context.Context, name, createdBy string) (Channel, error)
	ListChannels(ctx context.Context) ([]Channel, error)
	EnsureMember(ctx context.Context, channelID int64, userID string) error
	SaveChannelMessage(ctx context.Context, channelID int64, userID string, payload []byte) (Message, error)
	ListMessages(ctx context.Context, channelID int64, limit int) ([]Message, error)
	SaveMessage(ctx context.Context, subject string, payload []byte) error
	Close() error
}

// Presence tracks active connections.
type Presence interface {
	Incr(ctx context.Context, key string) error
	Decr(ctx context.Context, key string) error
	Close() error
}

// AuthConfig controls JWT auth.
type AuthConfig struct {
	Secret        []byte
	RefreshSecret []byte
	Enabled       bool
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	CookieDomain  string
	CookieSecure  bool
	CorsOrigin    string
}

// Channel model.
type Channel struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// Message model.
type Message struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	UserID    string    `json:"user_id"`
	Subject   string    `json:"subject"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// User model.
type User struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// RefreshToken model.
type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
}

// NewNatsAdapter adapts a *nats.Conn to the NatsClient interface.
func NewNatsAdapter(conn *nats.Conn) NatsClient {
	return &natsAdapter{conn: conn}
}

type natsAdapter struct {
	conn *nats.Conn
}

func (n *natsAdapter) Publish(subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

func (n *natsAdapter) ChanSubscribe(subject string, ch chan *nats.Msg) (Subscription, error) {
	sub, err := n.conn.ChanSubscribe(subject, ch)
	if err != nil {
		return nil, err
	}
	return &natsSubscription{sub: sub}, nil
}

func (n *natsAdapter) IsConnected() bool {
	return n.conn.IsConnected()
}

type natsSubscription struct {
	sub *nats.Subscription
}

func (s *natsSubscription) Unsubscribe() error {
	return s.sub.Unsubscribe()
}

func NewRouter(nc NatsClient, store Store, presence Presence, auth AuthConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(corsMiddleware(auth.CorsOrigin))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/ping-nats", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "nats not connected", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("nats ok"))
	})

	r.Route("/auth", func(ar chi.Router) {
		ar.Post("/register", func(w http.ResponseWriter, req *http.Request) {
			if store == nil {
				http.Error(w, "store not configured", http.StatusServiceUnavailable)
				return
			}
			var payload struct {
				UserID      string `json:"user_id"`
				Password    string `json:"password"`
				DisplayName string `json:"display_name"`
			}
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid payload", http.StatusBadRequest)
				return
			}
			payload.UserID = strings.TrimSpace(payload.UserID)
			payload.Password = strings.TrimSpace(payload.Password)
			if payload.UserID == "" || payload.Password == "" {
				http.Error(w, "user_id and password required", http.StatusBadRequest)
				return
			}
			user, err := store.CreateUser(req.Context(), payload.UserID, payload.Password, payload.DisplayName)
			if err != nil {
				log.Printf("create user failed: %v", err)
				http.Error(w, "create user failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusCreated, user)
		})

		ar.Post("/login", func(w http.ResponseWriter, req *http.Request) {
			if store == nil {
				http.Error(w, "store not configured", http.StatusServiceUnavailable)
				return
			}
			var payload struct {
				UserID   string `json:"user_id"`
				Password string `json:"password"`
			}
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid payload", http.StatusBadRequest)
				return
			}
			payload.UserID = strings.TrimSpace(payload.UserID)
			payload.Password = strings.TrimSpace(payload.Password)
			user, err := store.VerifyUserPassword(req.Context(), payload.UserID, payload.Password)
			if err != nil {
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
			issueSession(w, auth, store, user.ID)
			writeJSON(w, http.StatusOK, user)
		})

		ar.Post("/refresh", func(w http.ResponseWriter, req *http.Request) {
			if store == nil {
				http.Error(w, "store not configured", http.StatusServiceUnavailable)
				return
			}
			refreshToken := tokenFromCookie(req, "refresh_token")
			if refreshToken == "" {
				http.Error(w, "missing refresh token", http.StatusUnauthorized)
				return
			}
			claims := &jwt.RegisteredClaims{}
			parsed, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return auth.RefreshSecret, nil
			})
			if err != nil || !parsed.Valid || claims.Subject == "" {
				http.Error(w, "invalid refresh token", http.StatusUnauthorized)
				return
			}

			stored, err := store.GetRefreshToken(req.Context(), refreshToken)
			if err != nil || stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
				http.Error(w, "refresh token expired", http.StatusUnauthorized)
				return
			}

			_ = store.RevokeRefreshToken(req.Context(), refreshToken)
			issueSession(w, auth, store, claims.Subject)
			writeJSON(w, http.StatusOK, map[string]string{"status": "refreshed"})
		})

		ar.Post("/logout", func(w http.ResponseWriter, req *http.Request) {
			if store != nil {
				if token := tokenFromCookie(req, "refresh_token"); token != "" {
					_ = store.RevokeRefreshToken(req.Context(), token)
				}
			}
			clearSessionCookies(w, auth)
			writeJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
		})

		ar.With(authMiddleware(auth)).Get("/me", func(w http.ResponseWriter, req *http.Request) {
			userID := userFromContext(req.Context())
			if userID == "" {
				http.Error(w, "missing user", http.StatusUnauthorized)
				return
			}
			if store == nil {
				http.Error(w, "store not configured", http.StatusServiceUnavailable)
				return
			}
			user, err := store.GetUser(req.Context(), userID)
			if err != nil {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, user)
		})
	})

	r.Route("/", func(pr chi.Router) {
		pr.Use(authMiddleware(auth))

		pr.Post("/publish", func(w http.ResponseWriter, req *http.Request) {
			subject, err := subjectFromRequest(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			body, err := readBody(w, req)
			if err != nil {
				if errors.Is(err, errPayloadTooLarge) {
					http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(body) == 0 {
				body = []byte(`{"msg":"hello from gateway"}`)
			}

			if err := nc.Publish(subject, body); err != nil {
				http.Error(w, "publish failed: "+err.Error(), http.StatusBadGateway)
				return
			}
			if store != nil {
				if err := store.SaveMessage(req.Context(), subject, body); err != nil {
					log.Printf("store message failed: %v", err)
				}
			}
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte("published"))
		})

		pr.Get("/events", func(w http.ResponseWriter, req *http.Request) {
			subject, err := subjectFromRequest(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			streamSSE(w, req, nc, subject)
		})

		pr.Get("/ws", wsHandler(nc, store, presence))

		pr.Route("/channels", func(cr chi.Router) {
			cr.Get("/", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				channels, err := store.ListChannels(req.Context())
				if err != nil {
					log.Printf("list channels failed: %v", err)
					http.Error(w, "list channels failed: "+err.Error(), http.StatusInternalServerError)
					return
				}
				writeJSON(w, http.StatusOK, channels)
			})

			cr.Post("/", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				userID := userFromContext(req.Context())
				if userID == "" {
					http.Error(w, "missing user", http.StatusUnauthorized)
					return
				}
				var payload struct {
					Name string `json:"name"`
				}
				if err := json.NewDecoder(req.Body).Decode(&payload); err != nil || strings.TrimSpace(payload.Name) == "" {
					http.Error(w, "invalid payload", http.StatusBadRequest)
					return
				}

				if err := store.EnsureUser(req.Context(), userID); err != nil {
					http.Error(w, "ensure user failed", http.StatusInternalServerError)
					return
				}
				channel, err := store.CreateChannel(req.Context(), payload.Name, userID)
				if err != nil {
					log.Printf("create channel failed: %v", err)
					http.Error(w, "create channel failed: "+err.Error(), http.StatusInternalServerError)
					return
				}
				if err := store.EnsureMember(req.Context(), channel.ID, userID); err != nil {
					log.Printf("ensure member failed: %v", err)
				}
				writeJSON(w, http.StatusCreated, channel)
			})

			cr.Route("/{id}", func(ir chi.Router) {
				ir.Post("/messages", func(w http.ResponseWriter, req *http.Request) {
					if store == nil {
						http.Error(w, "store not configured", http.StatusServiceUnavailable)
						return
					}
					channelID, err := parseID(chi.URLParam(req, "id"))
					if err != nil {
						http.Error(w, "invalid channel id", http.StatusBadRequest)
						return
					}
					userID := userFromContext(req.Context())
					if userID == "" {
						http.Error(w, "missing user", http.StatusUnauthorized)
						return
					}
					if err := store.EnsureUser(req.Context(), userID); err != nil {
						http.Error(w, "ensure user failed", http.StatusInternalServerError)
						return
					}
					if err := store.EnsureMember(req.Context(), channelID, userID); err != nil {
						log.Printf("ensure member failed: %v", err)
					}

					payload, err := readMessagePayload(req)
					if err != nil {
						if errors.Is(err, errPayloadTooLarge) {
							http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
							return
						}
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					msg, err := store.SaveChannelMessage(req.Context(), channelID, userID, payload)
					if err != nil {
						log.Printf("save message failed: %v", err)
						http.Error(w, "save message failed: "+err.Error(), http.StatusInternalServerError)
						return
					}
					if err := nc.Publish(msg.Subject, payload); err != nil {
						log.Printf("nats publish failed: %v", err)
					}
					writeJSON(w, http.StatusCreated, msg)
				})

				ir.Get("/messages", func(w http.ResponseWriter, req *http.Request) {
					if store == nil {
						http.Error(w, "store not configured", http.StatusServiceUnavailable)
						return
					}
					channelID, err := parseID(chi.URLParam(req, "id"))
					if err != nil {
						http.Error(w, "invalid channel id", http.StatusBadRequest)
						return
					}
					limit := clamp(envIntFromQuery(req, "limit", 50), 1, 200)
					items, err := store.ListMessages(req.Context(), channelID, limit)
					if err != nil {
						log.Printf("list messages failed: %v", err)
						http.Error(w, "list messages failed: "+err.Error(), http.StatusInternalServerError)
						return
					}
					writeJSON(w, http.StatusOK, items)
				})
			})
		})

		pr.Route("/users", func(ur chi.Router) {
			ur.Get("/", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				users, err := store.ListUsers(req.Context())
				if err != nil {
					http.Error(w, "list users failed", http.StatusInternalServerError)
					return
				}
				writeJSON(w, http.StatusOK, users)
			})

			ur.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				user, err := store.GetUser(req.Context(), chi.URLParam(req, "id"))
				if err != nil {
					http.Error(w, "user not found", http.StatusNotFound)
					return
				}
				writeJSON(w, http.StatusOK, user)
			})

			ur.Post("/", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				var payload struct {
					UserID      string `json:"user_id"`
					Password    string `json:"password"`
					DisplayName string `json:"display_name"`
				}
				if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
					http.Error(w, "invalid payload", http.StatusBadRequest)
					return
				}
				payload.UserID = strings.TrimSpace(payload.UserID)
				payload.Password = strings.TrimSpace(payload.Password)
				if payload.UserID == "" || payload.Password == "" {
					http.Error(w, "user_id and password required", http.StatusBadRequest)
					return
				}
				user, err := store.CreateUser(req.Context(), payload.UserID, payload.Password, payload.DisplayName)
				if err != nil {
					http.Error(w, "create user failed", http.StatusInternalServerError)
					return
				}
				writeJSON(w, http.StatusCreated, user)
			})

			ur.Patch("/{id}", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				targetID := chi.URLParam(req, "id")
				userID := userFromContext(req.Context())
				if userID == "" || userID != targetID {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				var payload struct {
					DisplayName string `json:"display_name"`
					Password    string `json:"password"`
				}
				if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
					http.Error(w, "invalid payload", http.StatusBadRequest)
					return
				}
				user, err := store.UpdateUser(req.Context(), targetID, payload.DisplayName, payload.Password)
				if err != nil {
					http.Error(w, "update user failed", http.StatusInternalServerError)
					return
				}
				writeJSON(w, http.StatusOK, user)
			})

			ur.Delete("/{id}", func(w http.ResponseWriter, req *http.Request) {
				if store == nil {
					http.Error(w, "store not configured", http.StatusServiceUnavailable)
					return
				}
				targetID := chi.URLParam(req, "id")
				userID := userFromContext(req.Context())
				if userID == "" || userID != targetID {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				if err := store.DeleteUser(req.Context(), targetID); err != nil {
					http.Error(w, "delete user failed", http.StatusInternalServerError)
					return
				}
				writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
			})
		})
	})

	r.Handle("/metrics", promhttp.Handler())

	return r
}

func wsHandler(nc NatsClient, store Store, presence Presence) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool { return true },
	}

	return func(w http.ResponseWriter, req *http.Request) {
		subject, channelID, err := subjectOrChannel(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		userID := userFromContext(req.Context())
		if store != nil && channelID != 0 && userID != "" {
			if err := store.EnsureUser(ctx, userID); err != nil {
				log.Printf("ensure user failed: %v", err)
			}
			if err := store.EnsureMember(ctx, channelID, userID); err != nil {
				log.Printf("ensure member failed: %v", err)
			}
		}

		if presence != nil && channelID != 0 {
			if err := presence.Incr(ctx, presenceKey(channelID)); err != nil {
				log.Printf("presence incr failed: %v", err)
			}
			defer func() {
				if err := presence.Decr(context.Background(), presenceKey(channelID)); err != nil {
					log.Printf("presence decr failed: %v", err)
				}
			}()
		}

		ch := make(chan *nats.Msg, 256)
		sub, err := nc.ChanSubscribe(subject, ch)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("subscribe failed"))
			return
		}
		defer func() {
			_ = sub.Unsubscribe()
			close(ch)
		}()

		conn.SetReadLimit(maxBodyBytes)
		conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(90 * time.Second))
			return nil
		})

		done := make(chan struct{})
		pingTicker := time.NewTicker(30 * time.Second)
		go func() {
			defer close(done)
			defer pingTicker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-pingTicker.C:
					_ = conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second))
				case msg, ok := <-ch:
					if !ok {
						return
					}
					conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
					if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
						return
					}
				}
			}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				cancel()
				break
			}
			if len(message) == 0 {
				continue
			}
			if err := nc.Publish(subject, message); err != nil {
				log.Printf("ws publish failed: %v", err)
			}
			if store != nil && channelID != 0 && userID != "" {
				if _, err := store.SaveChannelMessage(ctx, channelID, userID, message); err != nil {
					log.Printf("store message failed: %v", err)
				}
			}
		}

		<-done
	}
}

func authMiddleware(cfg AuthConfig) func(http.Handler) http.Handler {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := tokenFromRequest(req)
			if token == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			claims := &jwt.RegisteredClaims{}
			parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return cfg.Secret, nil
			})
			if err != nil || !parsed.Valid || claims.Subject == "" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), ctxUserIDKey{}, claims.Subject)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

type ctxUserIDKey struct{}

func userFromContext(ctx context.Context) string {
	value := ctx.Value(ctxUserIDKey{})
	if value == nil {
		return ""
	}
	userID, _ := value.(string)
	return userID
}

func tokenFromRequest(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	if cookie := tokenFromCookie(req, "access_token"); cookie != "" {
		return cookie
	}
	return req.URL.Query().Get("token")
}

func tokenFromCookie(req *http.Request, name string) string {
	cookie, err := req.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func issueSession(w http.ResponseWriter, cfg AuthConfig, store Store, userID string) {
	accessToken, accessExp := signToken(cfg.Secret, userID, cfg.AccessTTL)
	refreshToken, refreshExp := signToken(cfg.RefreshSecret, userID, cfg.RefreshTTL)

	if store != nil {
		if err := store.SaveRefreshToken(context.Background(), userID, refreshToken, refreshExp); err != nil {
			log.Printf("save refresh token failed: %v", err)
		}
	}

	setCookie(w, "access_token", accessToken, accessExp, cfg)
	setCookie(w, "refresh_token", refreshToken, refreshExp, cfg)
}

func signToken(secret []byte, userID string, ttl time.Duration) (string, time.Time) {
	exp := time.Now().Add(ttl)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now()
	}
	return signed, exp
}

func setCookie(w http.ResponseWriter, name, value string, exp time.Time, cfg AuthConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  exp,
		MaxAge:   int(time.Until(exp).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cfg.CookieSecure,
		Domain:   cfg.CookieDomain,
	})
}

func clearSessionCookies(w http.ResponseWriter, cfg AuthConfig) {
	expired := time.Now().Add(-time.Hour)
	setCookie(w, "access_token", "", expired, cfg)
	setCookie(w, "refresh_token", "", expired, cfg)
}

func subjectFromRequest(req *http.Request) (string, error) {
	subject := req.URL.Query().Get("subject")
	if subject == "" {
		subject = defaultSubject
	}
	if !subjectRe.MatchString(subject) {
		return "", errors.New("invalid subject")
	}
	return subject, nil
}

func subjectOrChannel(req *http.Request) (string, int64, error) {
	if raw := req.URL.Query().Get("channel_id"); raw != "" {
		id, err := parseID(raw)
		if err != nil {
			return "", 0, errors.New("invalid channel id")
		}
		return channelSubject(id), id, nil
	}
	if raw := req.URL.Query().Get("subject"); raw != "" {
		if !subjectRe.MatchString(raw) {
			return "", 0, errors.New("invalid subject")
		}
		return raw, 0, nil
	}
	return defaultSubject, 0, nil
}

func channelSubject(id int64) string {
	return "channels." + strconv.FormatInt(id, 10)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envIntFromQuery(req *http.Request, key string, fallback int) int {
	if v := req.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		val, err := strconv.ParseBool(v)
		if err == nil {
			return val
		}
	}
	return fallback
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func writeSSE(w io.Writer, data string) error {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if _, err := io.WriteString(w, "data: "+line+"\n"); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, "\n")
	return err
}

func streamSSE(w http.ResponseWriter, req *http.Request, nc NatsClient, subject string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	_, _ = w.Write([]byte(": stream ready\n\n"))
	flusher.Flush()

	ch := make(chan *nats.Msg, 128)
	sub, err := nc.ChanSubscribe(subject, ch)
	if err != nil {
		http.Error(w, "subscribe failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer func() {
		_ = sub.Unsubscribe()
		close(ch)
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-req.Context().Done():
			return
		case <-ticker.C:
			_, _ = w.Write([]byte(": heartbeat\n\n"))
			flusher.Flush()
		case msg := <-ch:
			if msg == nil {
				continue
			}
			if err := writeSSE(w, string(msg.Data)); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func readBody(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	if w != nil {
		req.Body = http.MaxBytesReader(w, req.Body, maxBodyBytes)
	} else {
		req.Body = io.NopCloser(io.LimitReader(req.Body, maxBodyBytes+1))
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		if errors.Is(err, http.ErrBodyReadAfterClose) || strings.Contains(err.Error(), "http: request body too large") {
			return nil, errPayloadTooLarge
		}
		return nil, errors.New("cannot read body")
	}
	if len(body) > maxBodyBytes {
		return nil, errPayloadTooLarge
	}
	return body, nil
}

func readMessagePayload(req *http.Request) ([]byte, error) {
	body, err := readBody(nil, req)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errors.New("empty payload")
	}
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		var payload struct {
			Payload string `json:"payload"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, errors.New("invalid json payload")
		}
		if strings.TrimSpace(payload.Payload) == "" {
			return nil, errors.New("empty payload")
		}
		return []byte(payload.Payload), nil
	}
	return body, nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func presenceKey(channelID int64) string {
	return "channel:" + strconv.FormatInt(channelID, 10)
}

func corsMiddleware(origin string) func(http.Handler) http.Handler {
	if origin == "" {
		origin = "http://localhost:5173"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

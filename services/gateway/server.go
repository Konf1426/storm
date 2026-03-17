package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

var testMode = false

var (
	metricActiveWebSockets = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "storm_active_websockets",
		Help: "The total number of active WebSocket connections",
	})
	metricAuthRateLimited = promauto.NewCounter(prometheus.CounterOpts{
		Name: "storm_auth_rate_limited_total",
		Help: "The total number of auth requests rejected due to rate limiting",
	})
	metricNatsPublishDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "storm_nats_publish_duration_seconds",
		Help:    "Histogram of NATS publish latencies",
		Buckets: prometheus.DefBuckets,
	})
	metricSaveQueueLen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "storm_save_queue_length",
		Help: "The current number of messages waiting to be saved to database",
	})
)

type asyncTaskType int

const (
	taskSaveMessage asyncTaskType = iota
	taskSaveRefreshToken
)

type asyncTask struct {
	taskType  asyncTaskType
	channelID int64
	userID    string
	payload   []byte
	token     string
	expiresAt time.Time
}

var (
	asyncTaskQueue = make(chan asyncTask, 50000)
)

const (
	errStoreNotConfigured = "store not configured"
	errInvalidPayload     = "invalid payload"
	errUpdateUserFailed   = "update user failed"
	errMissingUser        = "missing user"
	errInvalidChannelID   = "invalid channel id"

	logEnsureMemberFailed = "ensure member failed: %v"
)

var noopCleanup = func() {
	// No presence tracking was enabled for this connection.
}

func StartWorkerPool(ctx context.Context, store Store, numWorkers int) {
	log.Printf("starting message worker pool with %d workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-asyncTaskQueue:
					metricSaveQueueLen.Set(float64(len(asyncTaskQueue)))
					processAsyncTask(id, store, task)
				}
			}
		}(i)
	}
}

func processAsyncTask(workerID int, store Store, task asyncTask) {
	switch task.taskType {
	case taskSaveMessage:
		if _, err := store.SaveChannelMessage(context.Background(), task.channelID, task.userID, task.payload); err != nil {
			log.Printf("worker %d: store message failed: %v", workerID, err)
		}
	case taskSaveRefreshToken:
		if err := store.SaveRefreshToken(context.Background(), task.userID, task.token, task.expiresAt); err != nil {
			log.Printf("worker %d: store refresh token failed: %v", workerID, err)
		}
	}
}

const (
	defaultSubject = "storm.events"
	maxBodyBytes   = 1 << 20
)

var subjectRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
var errPayloadTooLarge = errors.New("payload too large")

// rateLimiter provides per-IP token bucket rate limiting.
type rateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rateLimiterEntry
	r        rate.Limit
	b        int
}

type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newRateLimiter(r rate.Limit, b int) *rateLimiter {
	rl := &rateLimiter{
		limiters: make(map[string]*rateLimiterEntry),
		r:        r,
		b:        b,
	}
	go func() {
		for range time.Tick(5 * time.Minute) {
			rl.mu.Lock()
			for ip, e := range rl.limiters {
				if time.Since(e.lastSeen) > 10*time.Minute {
					delete(rl.limiters, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

func (rl *rateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	e, ok := rl.limiters[ip]
	if !ok {
		e = &rateLimiterEntry{limiter: rate.NewLimiter(rl.r, rl.b)}
		rl.limiters[ip] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

func (rl *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !envBool("AUTH_RATE_LIMIT_ENABLED", true) {
			next.ServeHTTP(w, req)
			return
		}

		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			ip = req.RemoteAddr
		}
		if !rl.get(ip).Allow() {
			metricAuthRateLimited.Inc()
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, req)
	})
}

// authRateLimiter: 5 req/min per IP, burst of 3.
var authRateLimiter = newRateLimiter(rate.Every(12*time.Second), 3)

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
	Secret        []byte // #nosec G117
	RefreshSecret []byte // #nosec G117
	Enabled       bool
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	CookieDomain  string
	CookieSecure  bool
	CorsOrigin    string
}

type tokenClaims struct {
	jwt.RegisteredClaims
	NoSnif     bool   `json:"nosnif"`
	URLDomaine string `json:"urldomaine,omitempty"`
}

type channelEventPayload struct {
	Type            string `json:"type"`
	MessageID       int64  `json:"message_id,omitempty"`
	ChannelID       int64  `json:"channel_id,omitempty"`
	UserID          string `json:"user_id,omitempty"`
	User            string `json:"user,omitempty"`
	Message         string `json:"message,omitempty"`
	RecipientUserID string `json:"recipient_user_id,omitempty"`
	ReaderID        string `json:"reader_id,omitempty"`
	Reader          string `json:"reader,omitempty"`
	ReadAt          string `json:"read_at,omitempty"`
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
	start := time.Now()
	err := n.conn.Publish(subject, data)
	metricNatsPublishDuration.Observe(time.Since(start).Seconds())
	return err
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
	r.Use(securityHeadersMiddleware(auth.CorsOrigin))
	r.Use(corsMiddleware(auth.CorsOrigin))

	registerHealthRoutes(r, nc)
	registerAuthRoutes(r, store, auth)
	registerProtectedRoutes(r, nc, store, presence, auth)

	return r
}

func registerHealthRoutes(r chi.Router, nc NatsClient) {
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Handle("/metrics", promhttp.Handler())

	r.Get("/ping-nats", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "nats not connected", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("nats ok"))
	})
}

func registerAuthRoutes(r chi.Router, store Store, auth AuthConfig) {
	r.Route("/auth", func(ar chi.Router) {
		ar.Use(authRateLimiter.middleware)
		ar.Post("/register", handleRegister(store))
		ar.Post("/login", handleLogin(store, auth))
		ar.Post("/refresh", handleRefresh(store, auth))
		ar.Post("/logout", handleLogout(store, auth))
		ar.With(authMiddleware(auth)).Get("/me", handleAuthMeGet(store))
		ar.With(authMiddleware(auth)).Patch("/me", handleAuthMePatch(store, auth))
		ar.With(authMiddleware(auth)).Delete("/me", handleAuthMeDelete(store, auth))
	})
}

func registerProtectedRoutes(r chi.Router, nc NatsClient, store Store, presence Presence, auth AuthConfig) {
	r.Route("/", func(pr chi.Router) {
		pr.Use(authMiddleware(auth))
		pr.Post("/publish", handlePublish(nc, store))
		pr.Get("/ws", wsHandler(nc, store, presence))
		registerChannelRoutes(pr, nc, store)
		registerUserRoutes(pr, store)
	})
}

func registerChannelRoutes(r chi.Router, nc NatsClient, store Store) {
	r.Route("/channels", func(cr chi.Router) {
		cr.Get("/", handleChannelsList(store))
		cr.Post("/", handleChannelsCreate(store))
		cr.Route("/{id}", func(ir chi.Router) {
			ir.Post("/messages", handleChannelMessagesCreate(nc, store))
			ir.Get("/messages", handleChannelMessagesList(store))
			ir.Post("/messages/{messageId}/read", handleChannelMessageRead(nc, store))
		})
	})
}

func registerUserRoutes(r chi.Router, store Store) {
	r.Route("/users", func(ur chi.Router) {
		ur.Get("/", handleUsersList(store))
		ur.Get("/{id}", handleUsersGet(store))
		ur.Post("/", handleUsersCreate(store))
		ur.Patch("/{id}", handleUsersPatch(store))
		ur.Delete("/{id}", handleUsersDelete(store))
	})
}

func handleRegister(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		var payload struct {
			UserID      string `json:"user_id"`
			Password    string `json:"password"` // #nosec G117
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
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
	}
}

func handleLogin(store Store, auth AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		var payload struct {
			UserID   string `json:"user_id"`
			Password string `json:"password"` // #nosec G117
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
			return
		}
		payload.UserID = strings.TrimSpace(payload.UserID)
		payload.Password = strings.TrimSpace(payload.Password)
		user, err := store.VerifyUserPassword(req.Context(), payload.UserID, payload.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if payload.UserID == user.DisplayName && user.ID != user.DisplayName {
			user, err = store.UpdateUser(req.Context(), user.ID, user.DisplayName, "")
			if err != nil {
				http.Error(w, errUpdateUserFailed, http.StatusInternalServerError)
				return
			}
		}
		issueSession(w, auth, store, user.ID)
		writeJSON(w, http.StatusOK, user)
	}
}

func handleRefresh(store Store, auth AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		refreshToken := tokenFromCookie(req, "refresh_token")
		if refreshToken == "" {
			http.Error(w, "missing refresh token", http.StatusUnauthorized)
			return
		}
		claims := &tokenClaims{}
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
	}
}

func handleLogout(store Store, auth AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if store != nil {
			if token := tokenFromCookie(req, "refresh_token"); token != "" {
				_ = store.RevokeRefreshToken(req.Context(), token)
			}
		}
		clearSessionCookies(w, auth)
		writeJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
	}
}

func handleAuthMeGet(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := requireAuthenticatedUser(w, req)
		if !ok || !requireStore(w, store) {
			return
		}
		user, err := store.GetUser(req.Context(), userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleAuthMePatch(store Store, auth AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := requireAuthenticatedUser(w, req)
		if !ok || !requireStore(w, store) {
			return
		}
		var payload struct {
			DisplayName string `json:"display_name"`
			Password    string `json:"password"` // #nosec G117
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
			return
		}
		payload.DisplayName = strings.TrimSpace(payload.DisplayName)
		payload.Password = strings.TrimSpace(payload.Password)
		if payload.DisplayName == "" && payload.Password == "" {
			http.Error(w, "display_name or password required", http.StatusBadRequest)
			return
		}
		user, err := store.UpdateUser(req.Context(), userID, payload.DisplayName, payload.Password)
		if err != nil {
			http.Error(w, errUpdateUserFailed, http.StatusInternalServerError)
			return
		}
		if user.ID != userID {
			if token := tokenFromCookie(req, "refresh_token"); token != "" {
				_ = store.RevokeRefreshToken(req.Context(), token)
			}
			issueSession(w, auth, store, user.ID)
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleAuthMeDelete(store Store, auth AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, ok := requireAuthenticatedUser(w, req)
		if !ok || !requireStore(w, store) {
			return
		}
		if token := tokenFromCookie(req, "refresh_token"); token != "" {
			_ = store.RevokeRefreshToken(req.Context(), token)
		}
		if err := store.DeleteUser(req.Context(), userID); err != nil {
			http.Error(w, "delete user failed", http.StatusInternalServerError)
			return
		}
		clearSessionCookies(w, auth)
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func handlePublish(nc NatsClient, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		subject, err := subjectFromRequest(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		body, ok := readPublishBody(w, req)
		if !ok {
			return
		}
		if err := nc.Publish(subject, body); err != nil {
			http.Error(w, "publish failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		savePublishedMessage(req.Context(), store, subject, body)
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("published"))
	}
}

func handleChannelsList(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		channels, err := store.ListChannels(req.Context())
		if err != nil {
			log.Printf("list channels failed: %v", err)
			http.Error(w, "list channels failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, channels)
	}
}

func handleChannelsCreate(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		userID, ok := requireAuthenticatedUser(w, req)
		if !ok {
			return
		}
		var payload struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil || strings.TrimSpace(payload.Name) == "" {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
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
		logEnsureMember(req.Context(), store, channel.ID, userID)
		writeJSON(w, http.StatusCreated, channel)
	}
}

func handleChannelMessagesCreate(nc NatsClient, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		channelID, userID, ok := requireChannelMessageAccess(w, req, store)
		if !ok {
			return
		}
		payload, ok := readChannelMessagePayload(w, req)
		if !ok {
			return
		}
		msg, err := store.SaveChannelMessage(req.Context(), channelID, userID, payload)
		if err != nil {
			log.Printf("save message failed: %v", err)
			http.Error(w, "save message failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := publishChannelMessageEvent(nc, msg, payload); err != nil {
			log.Printf("nats publish failed: %v", err)
		}
		writeJSON(w, http.StatusCreated, msg)
	}
}

func handleChannelMessageRead(nc NatsClient, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		channelID, readerID, ok := requireChannelMessageAccess(w, req, store)
		if !ok {
			return
		}
		messageID, err := parseID(chi.URLParam(req, "messageId"))
		if err != nil {
			http.Error(w, "invalid message id", http.StatusBadRequest)
			return
		}
		var payload struct {
			RecipientUserID string `json:"recipient_user_id"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil || strings.TrimSpace(payload.RecipientUserID) == "" {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
			return
		}
		payload.RecipientUserID = strings.TrimSpace(payload.RecipientUserID)
		if payload.RecipientUserID == readerID {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		readerName, err := resolveReaderName(req.Context(), store, readerID)
		if err != nil {
			http.Error(w, "reader lookup failed", http.StatusInternalServerError)
			return
		}
		if err := publishReadReceipt(nc, channelID, messageID, payload.RecipientUserID, readerID, readerName); err != nil {
			http.Error(w, "publish receipt failed", http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleChannelMessagesList(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		channelID, err := parseChannelID(chi.URLParam(req, "id"))
		if err != nil {
			http.Error(w, errInvalidChannelID, http.StatusBadRequest)
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
	}
}

func handleUsersList(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		users, err := store.ListUsers(req.Context())
		if err != nil {
			http.Error(w, "list users failed", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, users)
	}
}

func handleUsersGet(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		user, err := store.GetUser(req.Context(), chi.URLParam(req, "id"))
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleUsersCreate(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		var payload struct {
			UserID      string `json:"user_id"`
			Password    string `json:"password"` // #nosec G117
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
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
	}
}

func handleUsersPatch(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		targetID := chi.URLParam(req, "id")
		userID, ok := requireAuthenticatedUser(w, req)
		if !ok {
			return
		}
		if userID != targetID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		var payload struct {
			DisplayName string `json:"display_name"`
			Password    string `json:"password"` // #nosec G117
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, errInvalidPayload, http.StatusBadRequest)
			return
		}
		user, err := store.UpdateUser(req.Context(), targetID, payload.DisplayName, payload.Password)
		if err != nil {
			http.Error(w, errUpdateUserFailed, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleUsersDelete(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !requireStore(w, store) {
			return
		}
		targetID := chi.URLParam(req, "id")
		userID, ok := requireAuthenticatedUser(w, req)
		if !ok {
			return
		}
		if userID != targetID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err := store.DeleteUser(req.Context(), targetID); err != nil {
			http.Error(w, "delete user failed", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
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
			log.Printf("ws dial: %v", err)
			return
		}

		metricActiveWebSockets.Inc()
		defer func() {
			metricActiveWebSockets.Dec()
			_ = conn.Close()
		}()

		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		userID := userFromContext(req.Context())
		ensureWSUserMembership(ctx, store, channelID, userID)
		defer trackWSPresence(ctx, presence, channelID)()

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

		if err := configureWSConnection(conn); err != nil {
			return
		}

		done := startWSWriter(ctx, conn, ch)
		readWSMessages(cancel, conn, wsMessageSink{
			nc:        nc,
			store:     store,
			subject:   subject,
			channelID: channelID,
			userID:    userID,
		})

		<-done
	}
}

func ensureWSUserMembership(ctx context.Context, store Store, channelID int64, userID string) {
	if store == nil || channelID == 0 || userID == "" {
		return
	}
	if err := store.EnsureUser(ctx, userID); err != nil {
		log.Printf("ensure user failed: %v", err)
	}
	logEnsureMember(ctx, store, channelID, userID)
}

func trackWSPresence(ctx context.Context, presence Presence, channelID int64) func() {
	if presence == nil || channelID == 0 {
		return noopCleanup
	}
	key := presenceKey(channelID)
	if err := presence.Incr(ctx, key); err != nil {
		log.Printf("presence incr failed: %v", err)
	}
	return func() {
		if err := presence.Decr(context.Background(), key); err != nil {
			log.Printf("presence decr failed: %v", err)
		}
	}
}

func configureWSConnection(conn *websocket.Conn) error {
	conn.SetReadLimit(maxBodyBytes)
	if err := conn.SetReadDeadline(time.Now().Add(90 * time.Second)); err != nil {
		log.Printf("ws read deadline failed: %v", err)
		return err
	}
	conn.SetPongHandler(func(string) error {
		if err := conn.SetReadDeadline(time.Now().Add(90 * time.Second)); err != nil {
			log.Printf("ws pong deadline failed: %v", err)
			return err
		}
		return nil
	})
	return nil
}

func startWSWriter(ctx context.Context, conn *websocket.Conn, ch <-chan *nats.Msg) <-chan struct{} {
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
				if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
					log.Printf("ws write deadline failed: %v", err)
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
					return
				}
			}
		}
	}()
	return done
}

type wsMessageSink struct {
	nc        NatsClient
	store     Store
	subject   string
	channelID int64
	userID    string
}

func readWSMessages(cancel context.CancelFunc, conn *websocket.Conn, sink wsMessageSink) {
	defer cancel()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if len(message) == 0 {
			continue
		}
		if err := sink.nc.Publish(sink.subject, message); err != nil {
			log.Printf("ws publish failed: %v", err)
		}
		queueWSMessageSave(sink.store, sink.channelID, sink.userID, message)
	}
}

func queueWSMessageSave(store Store, channelID int64, userID string, message []byte) {
	if store == nil || channelID == 0 || userID == "" {
		return
	}
	// We copy the message because the original slice might be reused by the websocket reader.
	msgCopy := make([]byte, len(message))
	copy(msgCopy, message)
	select {
	case asyncTaskQueue <- asyncTask{taskType: taskSaveMessage, channelID: channelID, userID: userID, payload: msgCopy}:
		metricSaveQueueLen.Set(float64(len(asyncTaskQueue)))
	default:
		log.Printf("async task queue full, dropping message from %s", sanitize(userID)) // #nosec G706
	}
}

func requireStore(w http.ResponseWriter, store Store) bool {
	if store == nil {
		http.Error(w, errStoreNotConfigured, http.StatusServiceUnavailable)
		return false
	}
	return true
}

func requireAuthenticatedUser(w http.ResponseWriter, req *http.Request) (string, bool) {
	userID := userFromContext(req.Context())
	if userID == "" {
		http.Error(w, errMissingUser, http.StatusUnauthorized)
		return "", false
	}
	return userID, true
}

func parseChannelID(raw string) (int64, error) {
	return parseID(raw)
}

func logEnsureMember(ctx context.Context, store Store, channelID int64, userID string) {
	if err := store.EnsureMember(ctx, channelID, userID); err != nil {
		log.Printf(logEnsureMemberFailed, err)
	}
}

func readPublishBody(w http.ResponseWriter, req *http.Request) ([]byte, bool) {
	body, err := readBody(w, req)
	if err != nil {
		writePayloadError(w, err)
		return nil, false
	}
	if len(body) == 0 {
		return []byte(`{"msg":"hello from gateway"}`), true
	}
	return body, true
}

func savePublishedMessage(ctx context.Context, store Store, subject string, body []byte) {
	if store == nil {
		return
	}
	if err := store.SaveMessage(ctx, subject, body); err != nil {
		log.Printf("store message failed: %v", err)
	}
}

func publishChannelMessageEvent(nc NatsClient, msg Message, payload []byte) error {
	eventPayload, err := buildChannelMessageEvent(msg, payload)
	if err != nil {
		return err
	}
	return nc.Publish(msg.Subject, eventPayload)
}

func buildChannelMessageEvent(msg Message, payload []byte) ([]byte, error) {
	userID, user, message := parseChannelMessagePayload(payload, msg.UserID)
	event := channelEventPayload{
		Type:      "chat_message",
		MessageID: msg.ID,
		ChannelID: msg.ChannelID,
		UserID:    userID,
		User:      user,
		Message:   message,
	}
	return json.Marshal(event)
}

func parseChannelMessagePayload(payload []byte, fallbackUserID string) (string, string, string) {
	var data struct {
		UserID  string `json:"user_id"`
		User    string `json:"user"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(payload, &data); err == nil && strings.TrimSpace(data.Message) != "" {
		userID := strings.TrimSpace(data.UserID)
		if userID == "" {
			userID = fallbackUserID
		}
		return userID, strings.TrimSpace(data.User), data.Message
	}
	return fallbackUserID, "", string(payload)
}

func resolveReaderName(ctx context.Context, store Store, readerID string) (string, error) {
	user, err := store.GetUser(ctx, readerID)
	if err != nil {
		return "", err
	}
	if user.DisplayName != "" {
		return user.DisplayName, nil
	}
	return user.ID, nil
}

func publishReadReceipt(
	nc NatsClient,
	channelID int64,
	messageID int64,
	recipientUserID, readerID, readerName string,
) error {
	payload, err := json.Marshal(channelEventPayload{
		Type:            "read_receipt",
		MessageID:       messageID,
		ChannelID:       channelID,
		RecipientUserID: recipientUserID,
		ReaderID:        readerID,
		Reader:          readerName,
		ReadAt:          time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return nc.Publish(channelSubject(channelID), payload)
}

func requireChannelMessageAccess(w http.ResponseWriter, req *http.Request, store Store) (int64, string, bool) {
	if !requireStore(w, store) {
		return 0, "", false
	}
	channelID, err := parseChannelID(chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, errInvalidChannelID, http.StatusBadRequest)
		return 0, "", false
	}
	userID, ok := requireAuthenticatedUser(w, req)
	if !ok {
		return 0, "", false
	}
	if err := store.EnsureUser(req.Context(), userID); err != nil {
		http.Error(w, "ensure user failed", http.StatusInternalServerError)
		return 0, "", false
	}
	logEnsureMember(req.Context(), store, channelID, userID)
	return channelID, userID, true
}

func readChannelMessagePayload(w http.ResponseWriter, req *http.Request) ([]byte, bool) {
	payload, err := readMessagePayload(req)
	if err != nil {
		writePayloadError(w, err)
		return nil, false
	}
	return payload, true
}

func writePayloadError(w http.ResponseWriter, err error) {
	if errors.Is(err, errPayloadTooLarge) {
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
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

			claims := &tokenClaims{}
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
	accessToken, accessExp := signToken(cfg, cfg.Secret, userID, cfg.AccessTTL)
	refreshToken, refreshExp := signToken(cfg, cfg.RefreshSecret, userID, cfg.RefreshTTL)

	if store != nil {
		if testMode {
			if err := store.SaveRefreshToken(context.Background(), userID, refreshToken, refreshExp); err != nil {
				log.Printf("test sync save refresh token failed: %v", err)
			}
		} else {
			select {
			case asyncTaskQueue <- asyncTask{
				taskType:  taskSaveRefreshToken,
				userID:    userID,
				token:     refreshToken,
				expiresAt: refreshExp,
			}:
				metricSaveQueueLen.Set(float64(len(asyncTaskQueue)))
			default:
				log.Printf("async task queue full, dropping refresh token for %s", sanitize(userID)) // #nosec G706
			}
		}
	}

	setCookie(w, "access_token", accessToken, accessExp, cfg)
	setCookie(w, "refresh_token", refreshToken, refreshExp, cfg)
}

func signToken(cfg AuthConfig, secret []byte, userID string, ttl time.Duration) (string, time.Time) {
	exp := time.Now().Add(ttl)
	urlDomaine := cfg.CookieDomain
	if urlDomaine == "" {
		urlDomaine = cfg.CorsOrigin
	}
	claims := tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		NoSnif:     true,
		URLDomaine: urlDomaine,
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
			return "", 0, errors.New(errInvalidChannelID)
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

func securityHeadersMiddleware(corsOrigin string) func(http.Handler) http.Handler {
	csp := buildCSP(corsOrigin)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("X-DNS-Prefetch-Control", "off")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
			w.Header().Set("Content-Security-Policy", csp)
			if isHTTPS(req) {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			if !strings.HasPrefix(req.URL.Path, "/ws") {
				w.Header().Set("Cache-Control", "no-store")
			}
			next.ServeHTTP(w, req)
		})
	}
}

func buildCSP(corsOrigin string) string {
	connectSrc := []string{"'self'"}
	if corsOrigin != "" {
		connectSrc = append(connectSrc, corsOrigin)
		if wsOrigin := websocketOrigin(corsOrigin); wsOrigin != "" {
			connectSrc = append(connectSrc, wsOrigin)
		}
	}
	return strings.Join([]string{
		"default-src 'none'",
		"base-uri 'none'",
		"frame-ancestors 'none'",
		"form-action 'self'",
		"connect-src " + strings.Join(connectSrc, " "),
		"img-src 'self' data:",
		"style-src 'self' 'unsafe-inline'",
		"script-src 'self'",
	}, "; ")
}

func websocketOrigin(origin string) string {
	if strings.HasPrefix(origin, "https://") {
		return "wss://" + strings.TrimPrefix(origin, "https://")
	}
	if strings.HasPrefix(origin, "http://") {
		return "ws://" + strings.TrimPrefix(origin, "http://")
	}
	return ""
}

func isHTTPS(req *http.Request) bool {
	if req.TLS != nil {
		return true
	}
	return strings.EqualFold(req.Header.Get("X-Forwarded-Proto"), "https")
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

func sanitize(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", "")
}

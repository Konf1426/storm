package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type pgxPool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

type postgresStore struct {
	pool pgxPool
}

func NewPostgresStore(ctx context.Context, dsn string) (Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	s := &postgresStore{pool: pool}
	if err := s.ensureSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return s, nil
}

func newPostgresStoreWithPool(pool pgxPool) *postgresStore {
	return &postgresStore{pool: pool}
}

func (s *postgresStore) ensureSchema(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  password_hash TEXT NOT NULL,
  display_name TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS channels (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_by TEXT NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS channel_members (
  channel_id BIGINT NOT NULL REFERENCES channels(id),
  user_id TEXT NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (channel_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  channel_id BIGINT NULL REFERENCES channels(id),
  user_id TEXT NULL REFERENCES users(id),
  subject TEXT NOT NULL,
  payload BYTEA NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  token TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE messages ADD COLUMN IF NOT EXISTS channel_id BIGINT NULL;
ALTER TABLE messages ADD COLUMN IF NOT EXISTS user_id TEXT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';
`)
	return err
}

func (s *postgresStore) EnsureUser(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `INSERT INTO users (id, password_hash) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, "")
	return err
}

func (s *postgresStore) CreateUser(ctx context.Context, userID, password, displayName string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	var user User
	err = s.pool.QueryRow(ctx, `
INSERT INTO users (id, password_hash, display_name)
VALUES ($1, $2, $3)
RETURNING id, display_name, created_at
`, userID, string(hash), displayName).Scan(&user.ID, &user.DisplayName, &user.CreatedAt)
	return user, err
}

func (s *postgresStore) GetUser(ctx context.Context, userID string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user User
	err := s.pool.QueryRow(ctx, `SELECT id, display_name, created_at FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.DisplayName, &user.CreatedAt)
	return user, err
}

func (s *postgresStore) ListUsers(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, `SELECT id, display_name, created_at FROM users ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.DisplayName, &user.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, user)
	}
	return out, rows.Err()
}

func (s *postgresStore) UpdateUser(ctx context.Context, userID, displayName, password string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if displayName != "" {
		_, err := s.pool.Exec(ctx, `UPDATE users SET display_name = $1 WHERE id = $2`, displayName, userID)
		if err != nil {
			return User{}, err
		}
	}
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}
		_, err = s.pool.Exec(ctx, `UPDATE users SET password_hash = $1 WHERE id = $2`, string(hash), userID)
		if err != nil {
			return User{}, err
		}
	}

	return s.GetUser(ctx, userID)
}

func (s *postgresStore) DeleteUser(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	return err
}

func (s *postgresStore) VerifyUserPassword(ctx context.Context, userID, password string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var hash string
	var user User
	err := s.pool.QueryRow(ctx, `SELECT id, password_hash, display_name, created_at FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &hash, &user.DisplayName, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return User{}, errors.New("invalid password")
	}
	return user, nil
}

func (s *postgresStore) SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
INSERT INTO refresh_tokens (token, user_id, expires_at, revoked)
VALUES ($1, $2, $3, false)
`, token, userID, expiresAt)
	return err
}

func (s *postgresStore) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var out RefreshToken
	err := s.pool.QueryRow(ctx, `SELECT token, user_id, expires_at, revoked FROM refresh_tokens WHERE token = $1`, token).
		Scan(&out.Token, &out.UserID, &out.ExpiresAt, &out.Revoked)
	return out, err
}

func (s *postgresStore) RevokeRefreshToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked = true WHERE token = $1`, token)
	return err
}

func (s *postgresStore) CreateChannel(ctx context.Context, name, createdBy string) (Channel, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var channel Channel
	err := s.pool.QueryRow(ctx, `
INSERT INTO channels (name, created_by)
VALUES ($1, $2)
RETURNING id, name, created_by, created_at
`, name, createdBy).Scan(&channel.ID, &channel.Name, &channel.CreatedBy, &channel.CreatedAt)
	return channel, err
}

func (s *postgresStore) ListChannels(ctx context.Context) ([]Channel, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, `SELECT id, name, created_by, created_at FROM channels ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Channel
	for rows.Next() {
		var c Channel
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *postgresStore) EnsureMember(ctx context.Context, channelID int64, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
INSERT INTO channel_members (channel_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, channelID, userID)
	return err
}

func (s *postgresStore) SaveChannelMessage(ctx context.Context, channelID int64, userID string, payload []byte) (Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	subject := channelSubject(channelID)
	var msg Message
	msg.Subject = subject

	var dbPayload []byte
	err := s.pool.QueryRow(ctx, `
INSERT INTO messages (channel_id, user_id, subject, payload)
VALUES ($1, $2, $3, $4)
RETURNING id, channel_id, user_id, subject, payload, created_at
`, channelID, userID, subject, payload).Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Subject, &dbPayload, &msg.CreatedAt)
	if err == nil {
		msg.Payload = string(dbPayload)
	}
	return msg, err
}

func (s *postgresStore) ListMessages(ctx context.Context, channelID int64, limit int) ([]Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, `
SELECT id, channel_id, user_id, subject, payload, created_at
FROM messages
WHERE channel_id = $1
ORDER BY id DESC
LIMIT $2
`, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Message
	for rows.Next() {
		var msg Message
		var payload []byte
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Subject, &payload, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msg.Payload = string(payload)
		out = append(out, msg)
	}
	return out, rows.Err()
}

func (s *postgresStore) SaveMessage(ctx context.Context, subject string, payload []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, `INSERT INTO messages (subject, payload) VALUES ($1, $2)`, subject, payload)
	return err
}

func (s *postgresStore) Close() error {
	s.pool.Close()
	return nil
}

type redisPresence struct {
	client *redis.Client
}

func NewRedisPresence(ctx context.Context, addr, password string, db int) (Presence, error) {
	//nolint:gosec // Redis TLS is configured at infra level; local dev uses plaintext.
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &redisPresence{client: client}, nil
}

func (r *redisPresence) key(subject string) string {
	return fmt.Sprintf("presence:%s", subject)
}

func (r *redisPresence) Incr(ctx context.Context, key string) error {
	return r.client.Incr(ctx, r.key(key)).Err()
}

func (r *redisPresence) Decr(ctx context.Context, key string) error {
	return r.client.Decr(ctx, r.key(key)).Err()
}

func (r *redisPresence) Close() error {
	return r.client.Close()
}

func isPgNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

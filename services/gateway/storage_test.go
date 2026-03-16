package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"golang.org/x/crypto/bcrypt"
)

func TestPostgresStoreEnsureSchema(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	mock.ExpectExec("CREATE TABLE").WillReturnResult(pgxmock.NewResult("CREATE", 1))

	s := newPostgresStoreWithPool(mock)
	if err := s.ensureSchema(context.Background()); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestPostgresStoreUserFlow(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)

	mock.ExpectExec("INSERT INTO users").WithArgs("alice", pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	if err := s.EnsureUser(context.Background(), "alice"); err != nil {
		t.Fatalf("ensure user: %v", err)
	}

	mock.ExpectQuery("INSERT INTO users").WithArgs("bob", pgxmock.AnyArg(), "Bob").WillReturnRows(
		pgxmock.NewRows([]string{"id", "display_name", "created_at"}).AddRow("bob", "Bob", time.Now()),
	)
	if _, err := s.CreateUser(context.Background(), "bob", "pass", "Bob"); err != nil {
		t.Fatalf("create user: %v", err)
	}

	mock.ExpectQuery("SELECT id, display_name").WithArgs("bob").WillReturnRows(
		pgxmock.NewRows([]string{"id", "display_name", "created_at"}).AddRow("bob", "Bob", time.Now()),
	)
	if _, err := s.GetUser(context.Background(), "bob"); err != nil {
		t.Fatalf("get user: %v", err)
	}

	mock.ExpectQuery("SELECT id, display_name").WillReturnRows(
		pgxmock.NewRows([]string{"id", "display_name", "created_at"}).AddRow("bob", "Bob", time.Now()),
	)
	if _, err := s.ListUsers(context.Background()); err != nil {
		t.Fatalf("list users: %v", err)
	}

	mock.ExpectExec("UPDATE users SET display_name").WithArgs("Bob2", "bob").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("UPDATE users SET password_hash").WithArgs(pgxmock.AnyArg(), "bob").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectQuery("SELECT id, display_name").WithArgs("bob").WillReturnRows(
		pgxmock.NewRows([]string{"id", "display_name", "created_at"}).AddRow("bob", "Bob2", time.Now()),
	)
	if _, err := s.UpdateUser(context.Background(), "bob", "Bob2", "newpass"); err != nil {
		t.Fatalf("update user: %v", err)
	}

	mock.ExpectExec("DELETE FROM users").WithArgs("bob").WillReturnResult(pgxmock.NewResult("DELETE", 1))
	if err := s.DeleteUser(context.Background(), "bob"); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	mock.ExpectQuery("SELECT id, password_hash").WithArgs("eve").WillReturnRows(
		pgxmock.NewRows([]string{"id", "password_hash", "display_name", "created_at"}).AddRow("eve", string(hash), "Eve", time.Now()),
	)
	if _, err := s.VerifyUserPassword(context.Background(), "eve", "pass"); err != nil {
		t.Fatalf("verify user: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestPostgresStoreRefreshTokens(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)

	exp := time.Now().Add(1 * time.Hour)
	mock.ExpectExec("INSERT INTO refresh_tokens").WithArgs("token", "user", exp).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	if err := s.SaveRefreshToken(context.Background(), "user", "token", exp); err != nil {
		t.Fatalf("save refresh: %v", err)
	}

	mock.ExpectQuery("SELECT token, user_id").WithArgs("token").WillReturnRows(
		pgxmock.NewRows([]string{"token", "user_id", "expires_at", "revoked"}).AddRow("token", "user", exp, false),
	)
	if _, err := s.GetRefreshToken(context.Background(), "token"); err != nil {
		t.Fatalf("get refresh: %v", err)
	}

	mock.ExpectExec("UPDATE refresh_tokens").WithArgs("token").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	if err := s.RevokeRefreshToken(context.Background(), "token"); err != nil {
		t.Fatalf("revoke refresh: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestPostgresStoreChannelsAndMessages(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)

	mock.ExpectQuery("INSERT INTO channels").WithArgs("general", "alice").WillReturnRows(
		pgxmock.NewRows([]string{"id", "name", "created_by", "created_at"}).AddRow(int64(1), "general", "alice", time.Now()),
	)
	if _, err := s.CreateChannel(context.Background(), "general", "alice"); err != nil {
		t.Fatalf("create channel: %v", err)
	}

	mock.ExpectQuery("SELECT id, name").WillReturnRows(
		pgxmock.NewRows([]string{"id", "name", "created_by", "created_at"}).AddRow(int64(1), "general", "alice", time.Now()),
	)
	if _, err := s.ListChannels(context.Background()); err != nil {
		t.Fatalf("list channels: %v", err)
	}

	mock.ExpectExec("INSERT INTO channel_members").WithArgs(int64(1), "alice").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	if err := s.EnsureMember(context.Background(), 1, "alice"); err != nil {
		t.Fatalf("ensure member: %v", err)
	}

	payload := []byte("hello")
	mock.ExpectQuery("INSERT INTO messages").WithArgs(int64(1), "alice", "channels.1", payload).WillReturnRows(
		pgxmock.NewRows([]string{"id", "channel_id", "user_id", "subject", "payload", "created_at"}).AddRow(int64(1), int64(1), "alice", "channels.1", payload, time.Now()),
	)
	if _, err := s.SaveChannelMessage(context.Background(), 1, "alice", payload); err != nil {
		t.Fatalf("save channel message: %v", err)
	}

	mock.ExpectQuery("SELECT id, channel_id").WithArgs(int64(1), 10).WillReturnRows(
		pgxmock.NewRows([]string{"id", "channel_id", "user_id", "subject", "payload", "created_at"}).AddRow(int64(1), int64(1), "alice", "channels.1", payload, time.Now()),
	)
	if _, err := s.ListMessages(context.Background(), 1, 10); err != nil {
		t.Fatalf("list messages: %v", err)
	}

	mock.ExpectExec("INSERT INTO messages").WithArgs("storm.events", payload).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	if err := s.SaveMessage(context.Background(), "storm.events", payload); err != nil {
		t.Fatalf("save message: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestIsPgNotFound(t *testing.T) {
	if !isPgNotFound(pgx.ErrNoRows) {
		t.Fatalf("expected true")
	}
	if isPgNotFound(context.Canceled) {
		t.Fatalf("expected false")
	}
}

func TestPostgresStoreEnsureSchemaError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("boom"))
	if err := s.ensureSchema(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreVerifyPasswordInvalid(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT id, password_hash").WithArgs("bob").WillReturnRows(
		pgxmock.NewRows([]string{"id", "password_hash", "display_name", "created_at"}).AddRow("bob", "bad-hash", "Bob", time.Now()),
	)
	if _, err := s.VerifyUserPassword(context.Background(), "bob", "pass"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreUpdateUserNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("UPDATE users SET display_name").WithArgs("Bob2", "bob").WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	mock.ExpectQuery("SELECT id, display_name").WithArgs("bob").WillReturnError(pgx.ErrNoRows)
	if _, err := s.UpdateUser(context.Background(), "bob", "Bob2", ""); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreGetRefreshTokenNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT token, user_id").WithArgs("missing").WillReturnError(pgx.ErrNoRows)
	if _, err := s.GetRefreshToken(context.Background(), "missing"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreExecError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("INSERT INTO channel_members").WithArgs(int64(1), "alice").WillReturnError(errors.New("boom"))
	if err := s.EnsureMember(context.Background(), 1, "alice"); err == nil {
		t.Fatalf("expected error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestPostgresStoreListMessagesError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT id, channel_id").WithArgs(int64(1), 10).WillReturnError(errors.New("boom"))
	if _, err := s.ListMessages(context.Background(), 1, 10); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreSaveMessageError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("INSERT INTO messages").WithArgs("storm.events", []byte("x")).WillReturnError(errors.New("boom"))
	if err := s.SaveMessage(context.Background(), "storm.events", []byte("x")); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreGetUserNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT id, display_name").WithArgs("missing").WillReturnError(pgx.ErrNoRows)
	if _, err := s.GetUser(context.Background(), "missing"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreListUsersError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT id, display_name").WillReturnError(errors.New("boom"))
	if _, err := s.ListUsers(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreCreateChannelError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("INSERT INTO channels").WithArgs("general", "alice").WillReturnError(errors.New("boom"))
	if _, err := s.CreateChannel(context.Background(), "general", "alice"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreSaveChannelMessageError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("INSERT INTO messages").WithArgs(int64(1), "alice", "channels.1", []byte("x")).WillReturnError(errors.New("boom"))
	if _, err := s.SaveChannelMessage(context.Background(), 1, "alice", []byte("x")); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreCreateUserError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("INSERT INTO users").WithArgs("alice", pgxmock.AnyArg(), "Alice").WillReturnError(errors.New("boom"))
	if _, err := s.CreateUser(context.Background(), "alice", "pass", "Alice"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreSaveRefreshTokenError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("INSERT INTO refresh_tokens").WithArgs("token", "user", pgxmock.AnyArg()).WillReturnError(errors.New("boom"))
	if err := s.SaveRefreshToken(context.Background(), "user", "token", time.Now()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreRevokeRefreshTokenError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("UPDATE refresh_tokens").WithArgs("token").WillReturnError(errors.New("boom"))
	if err := s.RevokeRefreshToken(context.Background(), "token"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreVerifyUserPasswordQueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT id, password_hash").WithArgs("bob").WillReturnError(errors.New("boom"))
	if _, err := s.VerifyUserPassword(context.Background(), "bob", "pass"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreListChannelsError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("SELECT id, name").WillReturnError(errors.New("boom"))
	if _, err := s.ListChannels(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreEnsureUserError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("INSERT INTO users").WithArgs("alice", "").WillReturnError(errors.New("boom"))
	if err := s.EnsureUser(context.Background(), "alice"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreDeleteUserError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("DELETE FROM users").WithArgs("alice").WillReturnError(errors.New("boom"))
	if err := s.DeleteUser(context.Background(), "alice"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreListUsersRowsError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	rows := pgxmock.NewRows([]string{"id", "display_name", "created_at"}).AddRow("a", "A", time.Now()).RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, display_name").WillReturnRows(rows)
	if _, err := s.ListUsers(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreListChannelsRowsError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	rows := pgxmock.NewRows([]string{"id", "name", "created_by", "created_at"}).AddRow(int64(1), "c", "u", time.Now()).RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, name").WillReturnRows(rows)
	if _, err := s.ListChannels(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreListMessagesRowsError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	rows := pgxmock.NewRows([]string{"id", "channel_id", "user_id", "subject", "payload", "created_at"}).AddRow(int64(1), int64(1), "u", "s", []byte("x"), time.Now()).RowError(0, errors.New("row error"))
	mock.ExpectQuery("SELECT id, channel_id").WillReturnRows(rows)
	if _, err := s.ListMessages(context.Background(), 1, 10); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreSaveChannelMessageScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectQuery("INSERT INTO messages").WithArgs(int64(1), "alice", "channels.1", []byte("x")).WillReturnError(errors.New("scan error"))
	if _, err := s.SaveChannelMessage(context.Background(), 1, "alice", []byte("x")); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreListMessagesScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	rows := pgxmock.NewRows([]string{"id", "channel_id", "user_id", "subject", "payload", "created_at"}).AddRow("bad", int64(1), "u", "s", []byte("x"), time.Now())
	mock.ExpectQuery("SELECT id, channel_id").WillReturnRows(rows)
	if _, err := s.ListMessages(context.Background(), 1, 10); err == nil {
		t.Fatalf("expected error")
	}
}

func TestPostgresStoreSaveMessageExecErrorTag(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	defer mock.Close()

	s := newPostgresStoreWithPool(mock)
	mock.ExpectExec("INSERT INTO messages").WithArgs("storm.events", []byte("x")).WillReturnResult(pgxmock.NewResult("INSERT", 0))
	if err := s.SaveMessage(context.Background(), "storm.events", []byte("x")); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPostgresStoreClose(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	s := newPostgresStoreWithPool(mock)
	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestRedisPresenceIncrDecrClose(t *testing.T) {
	ctx := context.Background()
	srv := miniredis.RunT(t)
	pres, err := NewRedisPresence(ctx, srv.Addr(), "", 0)
	if err != nil {
		t.Fatalf("presence: %v", err)
	}
	rp, ok := pres.(*redisPresence)
	if !ok {
		t.Fatalf("expected redisPresence")
	}
	if got := rp.key("room"); got != "presence:room" {
		t.Fatalf("unexpected key %q", got)
	}
	if err := pres.Incr(ctx, "room"); err != nil {
		t.Fatalf("incr: %v", err)
	}
	got, err := srv.Get("presence:room")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != "1" {
		t.Fatalf("expected 1, got %q", got)
	}
	if err := pres.Decr(ctx, "room"); err != nil {
		t.Fatalf("decr: %v", err)
	}
	got, err = srv.Get("presence:room")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != "0" {
		t.Fatalf("expected 0, got %q", got)
	}
	if err := pres.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestPgconnCommandTag(t *testing.T) {
	var _ pgconn.CommandTag
}

func TestGetBcryptCost(t *testing.T) {
	os.Setenv("BCRYPT_COST", "")
	if cost := getBcryptCost(); cost != bcrypt.DefaultCost {
		t.Fatalf("expected DefaultCost, got %d", cost)
	}

	os.Setenv("BCRYPT_COST", "invalid")
	if cost := getBcryptCost(); cost != bcrypt.DefaultCost {
		t.Fatalf("expected DefaultCost for invalid, got %d", cost)
	}

	os.Setenv("BCRYPT_COST", "4")
	if cost := getBcryptCost(); cost != 4 {
		t.Fatalf("expected 4, got %d", cost)
	}
	os.Setenv("BCRYPT_COST", "")
}

func TestMaybeSimulateDelay(t *testing.T) {
	os.Setenv("SIMULATE_DB_DELAY", "")
	maybeSimulateDelay() // Should return immediately

	os.Setenv("SIMULATE_DB_DELAY", "invalid")
	maybeSimulateDelay() // Should return immediately

	os.Setenv("SIMULATE_DB_DELAY", "1ms")
	start := time.Now()
	maybeSimulateDelay()
	if time.Since(start) < 1*time.Millisecond {
		t.Fatalf("expected delay to occur")
	}
	os.Setenv("SIMULATE_DB_DELAY", "")
}

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

func initSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			token VARCHAR(255) NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL,
			revoked BOOLEAN NOT NULL DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed to execute schema setup: %w", err)
		}
	}

	// Seed admin if none exists
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		_, _ = db.Exec(
			"INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			uuid.NewString(), "admin", "admin@devserver.local", string(hash), time.Now().UTC(), time.Now().UTC(),
		)
	}

	return nil
}

func (s *MySQLStore) GetUserByUsername(ctx context.Context, username string) (*User, string, error) {
	var user User
	var hash string
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Email, &hash, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &user, hash, nil
}

func (s *MySQLStore) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.Email, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &user, nil
}

func (s *MySQLStore) CreateSession(ctx context.Context, session *Session) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO sessions (id, user_id, token, expires_at, created_at, revoked) VALUES (?, ?, ?, ?, ?, ?)",
		session.ID, session.UserID, session.Token, session.ExpiresAt.Format("2006-01-02 15:04:05"), session.CreatedAt.Format("2006-01-02 15:04:05"), session.Revoked,
	)
	return err
}

func (s *MySQLStore) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	var session Session
	var ea, ca []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, user_id, token, expires_at, created_at, revoked FROM sessions WHERE id = ?", sessionID).
		Scan(&session.ID, &session.UserID, &session.Token, &ea, &ca, &session.Revoked)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionExpired
		}
		return nil, err
	}

	session.ExpiresAt, _ = time.Parse("2006-01-02 15:04:05", string(ea))
	session.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))

	return &session, nil
}

func (s *MySQLStore) UpdateSession(ctx context.Context, session *Session) error {
	_, err := s.db.ExecContext(ctx, "UPDATE sessions SET revoked = ? WHERE id = ?", session.Revoked, session.ID)
	return err
}

func (s *MySQLStore) RevokeAllSessionsForUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE sessions SET revoked = 1 WHERE user_id = ?", userID)
	return err
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

package settings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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
		`CREATE TABLE IF NOT EXISTS settings (
			id VARCHAR(36) PRIMARY KEY,
			scope VARCHAR(32) NOT NULL,
			owner_id VARCHAR(36) NOT NULL,
			setting_key VARCHAR(255) NOT NULL,
			setting_value TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE KEY uniq_scope_owner_key (scope, owner_id, setting_key)
		);`,
		`CREATE TABLE IF NOT EXISTS integrations (
			id VARCHAR(36) PRIMARY KEY,
			scope VARCHAR(32) NOT NULL,
			owner_id VARCHAR(36) NOT NULL,
			provider VARCHAR(255) NOT NULL,
			status VARCHAR(64) NOT NULL,
			config JSON NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			INDEX idx_integration_owner (scope, owner_id)
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed to execute settings schema setup: %w", err)
		}
	}

	return nil
}

func (s *MySQLStore) GetSetting(ctx context.Context, scope Scope, ownerID, key string) (*Setting, error) {
	var setting Setting
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, scope, owner_id, setting_key, setting_value, created_at, updated_at FROM settings WHERE scope = ? AND owner_id = ? AND setting_key = ?", string(scope), ownerID, key).
		Scan(&setting.ID, &setting.Scope, &setting.OwnerID, &setting.Key, &setting.Value, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSettingNotFound
		}
		return nil, err
	}

	setting.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	setting.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &setting, nil
}

func (s *MySQLStore) ListSettings(ctx context.Context, scope Scope, ownerID string) ([]*Setting, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, scope, owner_id, setting_key, setting_value, created_at, updated_at FROM settings WHERE scope = ? AND owner_id = ?", string(scope), ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []*Setting
	for rows.Next() {
		var setting Setting
		var ca, ua []uint8
		if err := rows.Scan(&setting.ID, &setting.Scope, &setting.OwnerID, &setting.Key, &setting.Value, &ca, &ua); err != nil {
			return nil, err
		}
		setting.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
		setting.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))
		settings = append(settings, &setting)
	}

	return settings, nil
}

func (s *MySQLStore) UpsertSetting(ctx context.Context, setting *Setting) error {
	now := time.Now().UTC()
	if setting.ID == "" {
		setting.ID = uuid.NewString()
		setting.CreatedAt = now
	}
	setting.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO settings (id, scope, owner_id, setting_key, setting_value, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_at = VALUES(updated_at)`,
		setting.ID, string(setting.Scope), setting.OwnerID, setting.Key, setting.Value, setting.CreatedAt.Format("2006-01-02 15:04:05"), setting.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) DeleteSetting(ctx context.Context, scope Scope, ownerID, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM settings WHERE scope = ? AND owner_id = ? AND setting_key = ?", string(scope), ownerID, key)
	return err
}

func (s *MySQLStore) GetIntegration(ctx context.Context, id string) (*Integration, error) {
	var integration Integration
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, scope, owner_id, provider, status, config, created_at, updated_at FROM integrations WHERE id = ?", id).
		Scan(&integration.ID, &integration.Scope, &integration.OwnerID, &integration.Provider, &integration.Status, &integration.Config, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIntegrationNotFound
		}
		return nil, err
	}

	integration.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	integration.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &integration, nil
}

func (s *MySQLStore) ListIntegrations(ctx context.Context, scope Scope, ownerID string) ([]*Integration, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, scope, owner_id, provider, status, config, created_at, updated_at FROM integrations WHERE scope = ? AND owner_id = ?", string(scope), ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var integrations []*Integration
	for rows.Next() {
		var integration Integration
		var ca, ua []uint8
		if err := rows.Scan(&integration.ID, &integration.Scope, &integration.OwnerID, &integration.Provider, &integration.Status, &integration.Config, &ca, &ua); err != nil {
			return nil, err
		}
		integration.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
		integration.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))
		integrations = append(integrations, &integration)
	}

	return integrations, nil
}

func (s *MySQLStore) UpsertIntegration(ctx context.Context, integration *Integration) error {
	now := time.Now().UTC()
	if integration.ID == "" {
		integration.ID = uuid.NewString()
		integration.CreatedAt = now
	}
	integration.UpdatedAt = now

	if integration.Config == "" {
		integration.Config = "{}"
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO integrations (id, scope, owner_id, provider, status, config, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE provider = VALUES(provider), status = VALUES(status), config = VALUES(config), updated_at = VALUES(updated_at)`,
		integration.ID, string(integration.Scope), integration.OwnerID, integration.Provider, integration.Status, integration.Config, integration.CreatedAt.Format("2006-01-02 15:04:05"), integration.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) DeleteIntegration(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM integrations WHERE id = ?", id)
	return err
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

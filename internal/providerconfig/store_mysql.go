package providerconfig

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
	q := `CREATE TABLE IF NOT EXISTS provider_configs (
		id VARCHAR(36) PRIMARY KEY,
		scope VARCHAR(32) NOT NULL,
		owner_id VARCHAR(36) NOT NULL,
		type VARCHAR(64) NOT NULL,
		name VARCHAR(255) NOT NULL,
		enabled BOOLEAN NOT NULL DEFAULT FALSE,
		base_url VARCHAR(1024) NOT NULL DEFAULT '',
		default_model VARCHAR(255) NOT NULL DEFAULT '',
		secret_id VARCHAR(36) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		INDEX idx_provider_owner (scope, owner_id)
	);`

	if _, err := db.Exec(q); err != nil {
		return fmt.Errorf("failed to execute providerconfig schema setup: %w", err)
	}
	return nil
}

func (s *MySQLStore) GetConfig(ctx context.Context, id string) (*ProviderConfig, error) {
	var cfg ProviderConfig
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, scope, owner_id, type, name, enabled, base_url, default_model, secret_id, created_at, updated_at FROM provider_configs WHERE id = ?", id).
		Scan(&cfg.ID, &cfg.Scope, &cfg.OwnerID, &cfg.Type, &cfg.Name, &cfg.Enabled, &cfg.BaseURL, &cfg.DefaultModel, &cfg.SecretID, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("provider config not found")
		}
		return nil, err
	}

	var parseErr error
	if cfg.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	if cfg.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &cfg, nil
}

func (s *MySQLStore) ListConfigs(ctx context.Context, scope, ownerID string) ([]*ProviderConfig, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, scope, owner_id, type, name, enabled, base_url, default_model, secret_id, created_at, updated_at FROM provider_configs WHERE scope = ? AND owner_id = ?", scope, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*ProviderConfig
	for rows.Next() {
		var cfg ProviderConfig
		var ca, ua []uint8
		if err := rows.Scan(&cfg.ID, &cfg.Scope, &cfg.OwnerID, &cfg.Type, &cfg.Name, &cfg.Enabled, &cfg.BaseURL, &cfg.DefaultModel, &cfg.SecretID, &ca, &ua); err != nil {
			return nil, err
		}
		var parseErr error
		if cfg.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
		}
		if cfg.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
		}
		configs = append(configs, &cfg)
	}
	return configs, nil
}

func (s *MySQLStore) UpsertConfig(ctx context.Context, cfg *ProviderConfig) error {
	now := time.Now().UTC()
	if cfg.ID == "" {
		cfg.ID = uuid.NewString()
		cfg.CreatedAt = now
	}
	cfg.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO provider_configs (id, scope, owner_id, type, name, enabled, base_url, default_model, secret_id, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE type = VALUES(type), name = VALUES(name), enabled = VALUES(enabled), base_url = VALUES(base_url), default_model = VALUES(default_model), secret_id = VALUES(secret_id), updated_at = VALUES(updated_at)`,
		cfg.ID, cfg.Scope, cfg.OwnerID, string(cfg.Type), cfg.Name, cfg.Enabled, cfg.BaseURL, cfg.DefaultModel, cfg.SecretID, cfg.CreatedAt.Format("2006-01-02 15:04:05"), cfg.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) DeleteConfig(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM provider_configs WHERE id = ?", id)
	return err
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

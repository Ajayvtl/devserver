package environments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/Ajayvtl/devserver/internal/domain/common"
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
		`CREATE TABLE IF NOT EXISTS environments (
			id VARCHAR(36) PRIMARY KEY,
			owner_id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(32) NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			INDEX idx_owner (owner_id)
		);`,
		`CREATE TABLE IF NOT EXISTS env_variables (
			id VARCHAR(36) PRIMARY KEY,
			env_id VARCHAR(36) NOT NULL,
			var_key VARCHAR(255) NOT NULL,
			var_value TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (env_id) REFERENCES environments(id) ON DELETE CASCADE,
			UNIQUE KEY uniq_env_var (env_id, var_key)
		);`,
		`CREATE TABLE IF NOT EXISTS env_secrets (
			id VARCHAR(36) PRIMARY KEY,
			env_id VARCHAR(36) NOT NULL,
			secret_key VARCHAR(255) NOT NULL,
			secret_value TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (env_id) REFERENCES environments(id) ON DELETE CASCADE,
			UNIQUE KEY uniq_env_secret (env_id, secret_key)
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed to execute environments schema setup: %w", err)
		}
	}

	return nil
}

func (s *MySQLStore) UpsertEnvironment(ctx context.Context, env *Environment) error {
	now := time.Now().UTC()
	if env.ID == "" {
		env.ID = uuid.NewString()
		env.CreatedAt = now
	}
	env.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO environments (id, owner_id, name, type, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE name = VALUES(name), type = VALUES(type), updated_at = VALUES(updated_at)`,
		env.ID, env.OwnerID, env.Name, string(env.Type), env.CreatedAt.Format("2006-01-02 15:04:05"), env.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) GetEnvironment(ctx context.Context, id string) (*Environment, error) {
	var env Environment
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, owner_id, name, type, created_at, updated_at FROM environments WHERE id = ?", id).
		Scan(&env.ID, &env.OwnerID, &env.Name, &env.Type, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEnvNotFound
		}
		return nil, err
	}

	var parseErr error
	if env.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	if env.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &env, nil
}

func (s *MySQLStore) ListEnvironments(ctx context.Context, ownerID string, params common.QueryParams) ([]*Environment, int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM environments WHERE owner_id = ?", ownerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortCol := "created_at"
	if params.Sort == "name" {
		sortCol = "name"
	}

	dir := "DESC"
	if params.Dir == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf("SELECT id, owner_id, name, type, created_at, updated_at FROM environments WHERE owner_id = ? ORDER BY %s %s LIMIT ? OFFSET ?", sortCol, dir)
	rows, err := s.db.QueryContext(ctx, query, ownerID, params.Limit(), params.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var envs []*Environment
	for rows.Next() {
		var env Environment
		var ca, ua []uint8
		if err := rows.Scan(&env.ID, &env.OwnerID, &env.Name, &env.Type, &ca, &ua); err != nil {
			return nil, 0, err
		}

		var parseErr error
		if env.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
			return nil, 0, fmt.Errorf("failed to parse created_at: %w", parseErr)
		}
		if env.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
			return nil, 0, fmt.Errorf("failed to parse updated_at: %w", parseErr)
		}
		envs = append(envs, &env)
	}
	return envs, total, nil
}

func (s *MySQLStore) DeleteEnvironment(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM environments WHERE id = ?", id)
	return err
}

func (s *MySQLStore) UpsertVariable(ctx context.Context, variable *Variable) error {
	now := time.Now().UTC()
	if variable.ID == "" {
		variable.ID = uuid.NewString()
		variable.CreatedAt = now
	}
	variable.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO env_variables (id, env_id, var_key, var_value, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE var_value = VALUES(var_value), updated_at = VALUES(updated_at)`,
		variable.ID, variable.EnvID, variable.Key, variable.Value, variable.CreatedAt.Format("2006-01-02 15:04:05"), variable.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) GetVariable(ctx context.Context, envID, key string) (*Variable, error) {
	var v Variable
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, env_id, var_key, var_value, created_at, updated_at FROM env_variables WHERE env_id = ? AND var_key = ?", envID, key).
		Scan(&v.ID, &v.EnvID, &v.Key, &v.Value, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVarNotFound
		}
		return nil, err
	}

	var parseErr error
	if v.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	if v.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &v, nil
}

func (s *MySQLStore) ListVariables(ctx context.Context, envID string, params common.QueryParams) ([]*Variable, int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM env_variables WHERE env_id = ?", envID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortCol := "created_at"
	if params.Sort == "key" {
		sortCol = "var_key"
	}

	dir := "DESC"
	if params.Dir == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf("SELECT id, env_id, var_key, var_value, created_at, updated_at FROM env_variables WHERE env_id = ? ORDER BY %s %s LIMIT ? OFFSET ?", sortCol, dir)
	rows, err := s.db.QueryContext(ctx, query, envID, params.Limit(), params.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vars []*Variable
	for rows.Next() {
		var v Variable
		var ca, ua []uint8
		if err := rows.Scan(&v.ID, &v.EnvID, &v.Key, &v.Value, &ca, &ua); err != nil {
			return nil, 0, err
		}

		var parseErr error
		if v.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
			return nil, 0, fmt.Errorf("failed to parse created_at: %w", parseErr)
		}
		if v.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
			return nil, 0, fmt.Errorf("failed to parse updated_at: %w", parseErr)
		}
		vars = append(vars, &v)
	}
	return vars, total, nil
}

func (s *MySQLStore) DeleteVariable(ctx context.Context, envID, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM env_variables WHERE env_id = ? AND var_key = ?", envID, key)
	return err
}

func (s *MySQLStore) UpsertSecret(ctx context.Context, secret *Secret) error {
	now := time.Now().UTC()
	if secret.ID == "" {
		secret.ID = uuid.NewString()
		secret.CreatedAt = now
	}
	secret.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO env_secrets (id, env_id, secret_key, secret_value, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE secret_value = VALUES(secret_value), updated_at = VALUES(updated_at)`,
		secret.ID, secret.EnvID, secret.Key, secret.Value, secret.CreatedAt.Format("2006-01-02 15:04:05"), secret.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) GetSecret(ctx context.Context, envID, key string) (*Secret, error) {
	var sec Secret
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, env_id, secret_key, secret_value, created_at, updated_at FROM env_secrets WHERE env_id = ? AND secret_key = ?", envID, key).
		Scan(&sec.ID, &sec.EnvID, &sec.Key, &sec.Value, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSecretNotFound
		}
		return nil, err
	}

	var parseErr error
	if sec.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	if sec.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &sec, nil
}

func (s *MySQLStore) ListSecrets(ctx context.Context, envID string, params common.QueryParams) ([]*SecretReference, int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM env_secrets WHERE env_id = ?", envID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortCol := "id"
	if params.Sort == "key" {
		sortCol = "secret_key"
	}

	dir := "DESC"
	if params.Dir == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf("SELECT id, secret_key FROM env_secrets WHERE env_id = ? ORDER BY %s %s LIMIT ? OFFSET ?", sortCol, dir)
	rows, err := s.db.QueryContext(ctx, query, envID, params.Limit(), params.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var secrets []*SecretReference
	for rows.Next() {
		var ref SecretReference
		if err := rows.Scan(&ref.ID, &ref.Key); err != nil {
			return nil, 0, err
		}
		secrets = append(secrets, &ref)
	}
	return secrets, total, nil
}

func (s *MySQLStore) ListRawSecrets(ctx context.Context, envID string) ([]*Secret, error) {
	// INTERNAL USE ONLY: returns raw ciphertexts required for resolution expansion.
	rows, err := s.db.QueryContext(ctx, "SELECT id, env_id, secret_key, secret_value, created_at, updated_at FROM env_secrets WHERE env_id = ?", envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*Secret
	for rows.Next() {
		var sec Secret
		var ca, ua []uint8
		if err := rows.Scan(&sec.ID, &sec.EnvID, &sec.Key, &sec.Value, &ca, &ua); err != nil {
			return nil, err
		}
		var parseErr error
		if sec.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
		}
		if sec.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
		}
		secrets = append(secrets, &sec)
	}
	return secrets, nil
}

func (s *MySQLStore) DeleteSecret(ctx context.Context, envID, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM env_secrets WHERE env_id = ? AND secret_key = ?", envID, key)
	return err
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

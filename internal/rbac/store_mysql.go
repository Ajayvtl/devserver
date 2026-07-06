package rbac

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
		`CREATE TABLE IF NOT EXISTS organizations (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS roles (
			id VARCHAR(36) PRIMARY KEY,
			org_id VARCHAR(36),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			permissions JSON NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
			INDEX idx_org_id (org_id)
		);`,
		`CREATE TABLE IF NOT EXISTS memberships (
			id VARCHAR(36) PRIMARY KEY,
			org_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			role_id VARCHAR(36) NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			UNIQUE KEY uniq_org_user (org_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS resource_policies (
			resource_id VARCHAR(36) NOT NULL,
			resource_type VARCHAR(255) NOT NULL,
			org_id VARCHAR(36) NOT NULL,
			PRIMARY KEY (resource_id, resource_type),
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
			INDEX idx_policy_org (org_id)
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed to execute RBAC schema setup: %w", err)
		}
	}

	return nil
}

func (s *MySQLStore) GetOrganization(ctx context.Context, orgID string) (*Organization, error) {
	var org Organization
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, name, slug, created_at, updated_at FROM organizations WHERE id = ?", orgID).
		Scan(&org.ID, &org.Name, &org.Slug, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, err
	}

	org.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	org.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &org, nil
}

func (s *MySQLStore) GetMembership(ctx context.Context, userID, orgID string) (*Membership, error) {
	var mem Membership
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, org_id, user_id, role_id, created_at, updated_at FROM memberships WHERE user_id = ? AND org_id = ?", userID, orgID).
		Scan(&mem.ID, &mem.OrgID, &mem.UserID, &mem.RoleID, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccessDenied
		}
		return nil, err
	}

	mem.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	mem.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &mem, nil
}

func (s *MySQLStore) GetRole(ctx context.Context, roleID string) (*Role, error) {
	var role Role
	var orgID sql.NullString
	var desc sql.NullString
	var permJSON []uint8
	var ca, ua []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, org_id, name, description, permissions, created_at, updated_at FROM roles WHERE id = ?", roleID).
		Scan(&role.ID, &orgID, &role.Name, &desc, &permJSON, &ca, &ua)

	if err != nil {
		return nil, err
	}

	if orgID.Valid {
		role.OrgID = orgID.String
	}
	if desc.Valid {
		role.Description = desc.String
	}
	if err := json.Unmarshal(permJSON, &role.Permissions); err != nil {
		return nil, err
	}

	role.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ca))
	role.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", string(ua))

	return &role, nil
}

func (s *MySQLStore) GetResourcePolicy(ctx context.Context, resourceID, resourceType string) (*ResourcePolicy, error) {
	var pol ResourcePolicy

	err := s.db.QueryRowContext(ctx, "SELECT resource_id, resource_type, org_id FROM resource_policies WHERE resource_id = ? AND resource_type = ?", resourceID, resourceType).
		Scan(&pol.ResourceID, &pol.ResourceType, &pol.OrgID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return &pol, nil
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

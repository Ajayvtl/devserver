package rbac

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			last_login_at DATETIME NULL,
			joined_at DATETIME NULL,
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

	// Alter memberships table to add new columns if they do not exist
	_, _ = db.Exec("ALTER TABLE memberships ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'active'")
	_, _ = db.Exec("ALTER TABLE memberships ADD COLUMN last_login_at DATETIME NULL")
	_, _ = db.Exec("ALTER TABLE memberships ADD COLUMN joined_at DATETIME NULL")

	// Seed default org and admin membership if none exists
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&count)
	if count == 0 {
		// First try to find the admin user seeded by auth package
		var adminID string
		err := db.QueryRow("SELECT id FROM users WHERE username = 'admin' LIMIT 1").Scan(&adminID)
		if err == nil && adminID != "" {
			orgID := uuid.NewString()
			roleID := uuid.NewString()
			now := time.Now().UTC().Format("2006-01-02 15:04:05")

			// Create Default Organization
			_, _ = db.Exec("INSERT INTO organizations (id, name, slug, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
				orgID, "DevServer Default", "devserver-default", now, now)

			// Create Admin Role
			_, _ = db.Exec("INSERT INTO roles (id, org_id, name, permissions, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
				roleID, orgID, "admin", `["*"]`, now, now)

			// Assign Admin to Default Organization
			_, _ = db.Exec("INSERT INTO memberships (id, org_id, user_id, role_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
				uuid.NewString(), orgID, adminID, roleID, now, now)
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

	var parseErr error
	org.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	org.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &org, nil
}

func (s *MySQLStore) GetMembership(ctx context.Context, userID, orgID string) (*Membership, error) {
	var mem Membership
	var ca, ua []uint8
	var lla, ja []uint8

	err := s.db.QueryRowContext(ctx, "SELECT id, org_id, user_id, role_id, status, last_login_at, joined_at, created_at, updated_at FROM memberships WHERE user_id = ? AND org_id = ?", userID, orgID).
		Scan(&mem.ID, &mem.OrgID, &mem.UserID, &mem.RoleID, &mem.Status, &lla, &ja, &ca, &ua)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccessDenied
		}
		return nil, err
	}

	if len(lla) > 0 {
		t, err := time.Parse("2006-01-02 15:04:05", string(lla))
		if err == nil {
			mem.LastLoginAt = &t
		}
	}
	if len(ja) > 0 {
		t, err := time.Parse("2006-01-02 15:04:05", string(ja))
		if err == nil {
			mem.JoinedAt = &t
		}
	}

	var parseErr error
	mem.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	mem.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

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

	var parseErr error
	role.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}
	role.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &role, nil
}

func (s *MySQLStore) UpsertOrganization(ctx context.Context, org *Organization) error {
	now := time.Now().UTC()
	if org.ID == "" {
		org.ID = uuid.NewString()
		org.CreatedAt = now
	}
	org.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO organizations (id, name, created_at, updated_at) 
		 VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE name = VALUES(name), updated_at = VALUES(updated_at)`,
		org.ID, org.Name, org.CreatedAt.Format("2006-01-02 15:04:05"), org.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) UpsertRole(ctx context.Context, role *Role) error {
	now := time.Now().UTC()
	if role.ID == "" {
		role.ID = uuid.NewString()
		role.CreatedAt = now
	}
	role.UpdatedAt = now

	permsJSON, _ := json.Marshal(role.Permissions)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO roles (id, org_id, name, permissions, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE name = VALUES(name), permissions = VALUES(permissions), updated_at = VALUES(updated_at)`,
		role.ID, role.OrgID, role.Name, string(permsJSON), role.CreatedAt.Format("2006-01-02 15:04:05"), role.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) UpsertMembership(ctx context.Context, mem *Membership) error {
	now := time.Now().UTC()
	if mem.ID == "" {
		mem.ID = uuid.NewString()
		mem.CreatedAt = now
	}
	mem.UpdatedAt = now
	if mem.Status == "" {
		mem.Status = "active"
	}

	var llaVal, jaVal any
	if mem.LastLoginAt != nil {
		llaVal = mem.LastLoginAt.Format("2006-01-02 15:04:05")
	}
	if mem.JoinedAt != nil {
		jaVal = mem.JoinedAt.Format("2006-01-02 15:04:05")
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO memberships (id, user_id, org_id, role_id, status, last_login_at, joined_at, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE role_id = VALUES(role_id), status = VALUES(status), last_login_at = VALUES(last_login_at), joined_at = VALUES(joined_at), updated_at = VALUES(updated_at)`,
		mem.ID, mem.UserID, mem.OrgID, mem.RoleID, mem.Status, llaVal, jaVal, mem.CreatedAt.Format("2006-01-02 15:04:05"), mem.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (s *MySQLStore) ProvisionOrganizationTx(ctx context.Context, org *Organization, role *Role, mem *Membership) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	if org.ID == "" {
		org.ID = uuid.NewString()
		org.CreatedAt = now
	}
	org.UpdatedAt = now

	slug := strings.ToLower(strings.ReplaceAll(org.Name, " ", "-"))
	_, err = tx.ExecContext(ctx,
		`INSERT INTO organizations (id, name, slug, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		org.ID, org.Name, slug, org.CreatedAt.Format("2006-01-02 15:04:05"), org.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	role.OrgID = org.ID
	if role.ID == "" {
		role.ID = uuid.NewString()
		role.CreatedAt = now
	}
	role.UpdatedAt = now

	permsJSON, _ := json.Marshal(role.Permissions)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO roles (id, org_id, name, permissions, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		role.ID, role.OrgID, role.Name, string(permsJSON), role.CreatedAt.Format("2006-01-02 15:04:05"), role.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	mem.OrgID = org.ID
	mem.RoleID = role.ID
	if mem.ID == "" {
		mem.ID = uuid.NewString()
		mem.CreatedAt = now
	}
	mem.UpdatedAt = now

	_, err = tx.ExecContext(ctx,
		`INSERT INTO memberships (id, user_id, org_id, role_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		mem.ID, mem.UserID, mem.OrgID, mem.RoleID, mem.CreatedAt.Format("2006-01-02 15:04:05"), mem.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *MySQLStore) ListOrganizationsForUser(ctx context.Context, userID string, params common.QueryParams) ([]*Organization, int, error) {
	// First get total count
	var total int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(o.id)
		FROM organizations o
		JOIN memberships m ON o.id = m.org_id
		WHERE m.user_id = ?
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Determine sort column to avoid SQL injection
	sortCol := "o.created_at"
	if params.Sort == "name" {
		sortCol = "o.name"
	} else if params.Sort == "updated_at" {
		sortCol = "o.updated_at"
	}

	dir := "DESC"
	if params.Dir == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT o.id, o.name, o.created_at, o.updated_at 
		FROM organizations o
		JOIN memberships m ON o.id = m.org_id
		WHERE m.user_id = ?
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, sortCol, dir)

	rows, err := s.db.QueryContext(ctx, query, userID, params.Limit(), params.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orgs []*Organization
	for rows.Next() {
		var org Organization
		var ca, ua []uint8
		if err := rows.Scan(&org.ID, &org.Name, &ca, &ua); err != nil {
			return nil, 0, err
		}
		var parseErr error
		if org.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ca)); parseErr != nil {
			return nil, 0, fmt.Errorf("failed to parse created_at: %w", parseErr)
		}
		if org.UpdatedAt, parseErr = time.Parse("2006-01-02 15:04:05", string(ua)); parseErr != nil {
			return nil, 0, fmt.Errorf("failed to parse updated_at: %w", parseErr)
		}
		orgs = append(orgs, &org)
	}
	return orgs, total, nil
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

func (s *MySQLStore) ListMemberships(ctx context.Context, orgID string, params common.QueryParams) ([]*MembershipDetails, int, error) {
	where := "WHERE m.org_id = ?"
	args := []any{orgID}
	if search := params.Query; search != "" {
		where += " AND (u.username LIKE ? OR u.email LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(m.id)
		FROM memberships m
		JOIN users u ON m.user_id = u.id
		%s
	`, where)
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortCol := "m.created_at"
	if params.Sort == "username" {
		sortCol = "u.username"
	} else if params.Sort == "email" {
		sortCol = "u.email"
	} else if params.Sort == "status" {
		sortCol = "m.status"
	}

	dir := "DESC"
	if params.Dir == "ASC" {
		dir = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT m.id, m.org_id, m.user_id, u.username, u.email, m.role_id, r.name, m.status, m.last_login_at, m.joined_at, m.created_at, m.updated_at
		FROM memberships m
		JOIN users u ON m.user_id = u.id
		LEFT JOIN roles r ON m.role_id = r.id
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, where, sortCol, dir)

	args = append(args, params.Limit(), params.Offset())
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var details []*MembershipDetails
	for rows.Next() {
		var det MembershipDetails
		var roleName sql.NullString
		var lla, ja []uint8
		var ca, ua []uint8

		err := rows.Scan(
			&det.ID, &det.OrgID, &det.UserID, &det.Username, &det.Email,
			&det.RoleID, &roleName, &det.Status, &lla, &ja, &ca, &ua,
		)
		if err != nil {
			return nil, 0, err
		}

		det.RoleName = roleName.String
		if len(lla) > 0 {
			t, err := time.Parse("2006-01-02 15:04:05", string(lla))
			if err == nil {
				det.LastLoginAt = &t
			}
		}
		if len(ja) > 0 {
			t, err := time.Parse("2006-01-02 15:04:05", string(ja))
			if err == nil {
				det.JoinedAt = &t
			}
		}

		if det.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(ca)); err != nil {
			return nil, 0, err
		}
		if det.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(ua)); err != nil {
			return nil, 0, err
		}

		details = append(details, &det)
	}

	return details, total, nil
}

func (s *MySQLStore) UpdateMembershipStatus(ctx context.Context, orgID, userID string, status string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE memberships SET status = ?, updated_at = NOW() WHERE org_id = ? AND user_id = ?",
		status, orgID, userID,
	)
	return err
}

func (s *MySQLStore) UpdateMembershipRole(ctx context.Context, orgID, userID string, roleID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE memberships SET role_id = ?, updated_at = NOW() WHERE org_id = ? AND user_id = ?",
		roleID, orgID, userID,
	)
	return err
}

func (s *MySQLStore) RemoveMembership(ctx context.Context, orgID, userID string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM memberships WHERE org_id = ? AND user_id = ?",
		orgID, userID,
	)
	return err
}

func (s *MySQLStore) InviteMember(ctx context.Context, orgID string, email string, roleID string) error {
	var userID string
	err := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userID = uuid.NewString()
			username := strings.Split(email, "@")[0]
			username = fmt.Sprintf("%s-%s", username, userID[:8])
			_, err = s.db.ExecContext(ctx,
				"INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, 'placeholder', NOW(), NOW())",
				userID, username, email,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	memID := uuid.NewString()
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO memberships (id, user_id, org_id, role_id, status, joined_at, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, 'invited', NULL, NOW(), NOW())
		 ON DUPLICATE KEY UPDATE status = 'invited', role_id = VALUES(role_id), updated_at = NOW()`,
		memID, userID, orgID, roleID,
	)
	return err
}

func (s *MySQLStore) ListRoles(ctx context.Context, orgID string) ([]*Role, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, org_id, name, permissions FROM roles WHERE org_id = ?", orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*Role
	for rows.Next() {
		var role Role
		var permsJSON []byte
		if err := rows.Scan(&role.ID, &role.OrgID, &role.Name, &permsJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(permsJSON, &role.Permissions); err != nil {
			role.Permissions = []string{}
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

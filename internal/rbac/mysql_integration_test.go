package rbac

import (
    "context"
    "os"
    "testing"

    "github.com/Ajayvtl/devserver/internal/domain/common"
)

// These integration tests require a running MySQL instance and a DSN provided
// via the TEST_MYSQL_DSN environment variable. They are skipped when unset.

func TestMySQLStore_UpsertAndGet(t *testing.T) {
    dsn := os.Getenv("TEST_MYSQL_DSN")
    if dsn == "" {
        t.Skip("TEST_MYSQL_DSN not set; skipping MySQL integration test")
    }

    s, err := NewMySQLStore(dsn)
    if err != nil {
        t.Fatalf("failed to create MySQL store: %v", err)
    }
    defer s.Close()

    ctx := context.Background()

    org := &Organization{ID: "test-org-1", Name: "Test Org"}
    if err := s.UpsertOrganization(ctx, org); err != nil {
        t.Fatalf("UpsertOrganization failed: %v", err)
    }

    gotOrg, err := s.GetOrganization(ctx, org.ID)
    if err != nil {
        t.Fatalf("GetOrganization failed: %v", err)
    }
    if gotOrg.Name != org.Name {
        t.Fatalf("expected org name %s, got %s", org.Name, gotOrg.Name)
    }

    role := &Role{ID: "test-role-1", OrgID: org.ID, Name: "tester", Permissions: []string{PermissionMembersView}}
    if err := s.UpsertRole(ctx, role); err != nil {
        t.Fatalf("UpsertRole failed: %v", err)
    }

    gotRole, err := s.GetRole(ctx, role.ID)
    if err != nil {
        t.Fatalf("GetRole failed: %v", err)
    }
    if gotRole.Name != role.Name {
        t.Fatalf("expected role name %s, got %s", role.Name, gotRole.Name)
    }

    mem := &Membership{ID: "test-mem-1", OrgID: org.ID, UserID: "user-test-1", RoleID: role.ID}
    if err := s.UpsertMembership(ctx, mem); err != nil {
        t.Fatalf("UpsertMembership failed: %v", err)
    }

    gotMem, err := s.GetMembership(ctx, mem.UserID, mem.OrgID)
    if err != nil {
        t.Fatalf("GetMembership failed: %v", err)
    }
    if gotMem.RoleID != role.ID {
        t.Fatalf("expected role id %s, got %s", role.ID, gotMem.RoleID)
    }
}

func TestMySQLStore_ProvisionOrganizationTx(t *testing.T) {
    dsn := os.Getenv("TEST_MYSQL_DSN")
    if dsn == "" {
        t.Skip("TEST_MYSQL_DSN not set; skipping MySQL integration test")
    }

    s, err := NewMySQLStore(dsn)
    if err != nil {
        t.Fatalf("failed to create MySQL store: %v", err)
    }
    defer s.Close()

    ctx := context.Background()

    org := &Organization{Name: "Provisioned Org"}
    role := &Role{Name: "owner", Permissions: []string{PermissionWildcard}}
    mem := &Membership{UserID: "prov-user-1"}

    if err := s.ProvisionOrganizationTx(ctx, org, role, mem); err != nil {
        t.Fatalf("ProvisionOrganizationTx failed: %v", err)
    }

    // basic assertions: org should now be retrievable via ListAllOrganizations
    orgs, _, err := s.ListAllOrganizations(ctx, common.QueryParams{})
    if err != nil {
        t.Fatalf("ListAllOrganizations failed: %v", err)
    }
    if len(orgs) == 0 {
        t.Fatalf("expected at least one organization after provisioning")
    }
}

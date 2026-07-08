package rbac

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
)

// This file contains an expanded set of server-side tests targeting
// Step 1 coverage requirements: tenant isolation, permission registry,
// revoked/deleted role/member behavior, permission refresh, invalid perms.

func TestPermissionRegistryDefined(t *testing.T) {
	// ensure some known permissions are present
	if PermissionMembersView == "" || PermissionMembersManage == "" || PermissionWildcard == "" {
		t.Fatalf("expected permission constants to be defined")
	}
}

func TestTenantIsolationAndCrossTenantDenial(t *testing.T) {
	ms := newMockStore()
	// user1 is member of orgA only
	r := &Role{ID: "r1", Name: "dev", Permissions: []string{PermissionMembersView}}
	ms.roles[r.ID] = r
	mem := &Membership{ID: "m1", OrgID: "orgA", UserID: "user1", RoleID: r.ID}
	ms.memberships[mem.ID] = mem

	svc := NewService(zerolog.Nop(), ms)

	// allowed in orgA
	if err := svc.Authorize(context.Background(), "user1", "orgA", PermissionMembersView); err != nil {
		t.Fatalf("expected allowed in same org, got %v", err)
	}

	// denied in orgB
	if err := svc.Authorize(context.Background(), "user1", "orgB", PermissionMembersView); err == nil {
		t.Fatalf("expected denied in other org")
	}
}

func TestRevokedPermissionAndDeletedRoleBehavior(t *testing.T) {
	ms := newMockStore()
	// role initially has view
	r := &Role{ID: "r2", Name: "role2", Permissions: []string{PermissionMembersView}}
	ms.roles[r.ID] = r
	mem := &Membership{ID: "m2", OrgID: "orgX", UserID: "u2", RoleID: r.ID}
	ms.memberships[mem.ID] = mem

	svc := NewService(zerolog.Nop(), ms)

	// allowed initially
	if err := svc.Authorize(context.Background(), "u2", "orgX", PermissionMembersView); err != nil {
		t.Fatalf("expected allowed initially, got %v", err)
	}

	// revoke permission by removing from role
	r.Permissions = []string{}
	ms.roles[r.ID] = r

	if err := svc.Authorize(context.Background(), "u2", "orgX", PermissionMembersView); err == nil {
		t.Fatalf("expected denied after revoke")
	}

	// deleted role behavior: remove role entry
	delete(ms.roles, r.ID)
	if err := svc.Authorize(context.Background(), "u2", "orgX", PermissionMembersView); err == nil {
		t.Fatalf("expected denied when role missing")
	}
}

func TestDeletedMembershipBehavior(t *testing.T) {
	ms := newMockStore()
	r := &Role{ID: "r3", Name: "viewer", Permissions: []string{PermissionMembersView}}
	ms.roles[r.ID] = r
	mem := &Membership{ID: "m3", OrgID: "orgZ", UserID: "u3", RoleID: r.ID}
	ms.memberships[mem.ID] = mem

	svc := NewService(zerolog.Nop(), ms)

	if err := svc.Authorize(context.Background(), "u3", "orgZ", PermissionMembersView); err != nil {
		t.Fatalf("expected allowed before removal")
	}

	// remove membership
	delete(ms.memberships, mem.ID)
	if err := svc.Authorize(context.Background(), "u3", "orgZ", PermissionMembersView); err == nil {
		t.Fatalf("expected denied after membership removed")
	}
}

func TestGetUserPermissions_OrganizationSwitchAndRefresh(t *testing.T) {
	ms := newMockStore()
	rA := &Role{ID: "rA", Name: "admin", Permissions: []string{PermissionMembersManage}}
	rB := &Role{ID: "rB", Name: "viewer", Permissions: []string{PermissionMembersView}}
	ms.roles[rA.ID] = rA
	ms.roles[rB.ID] = rB

	memA := &Membership{ID: "mA", OrgID: "org1", UserID: "u-switch", RoleID: rA.ID}
	memB := &Membership{ID: "mB", OrgID: "org2", UserID: "u-switch", RoleID: rB.ID}
	ms.memberships[memA.ID] = memA
	ms.memberships[memB.ID] = memB

	svc := NewService(zerolog.Nop(), ms)

	perms1, err := svc.GetUserPermissions(context.Background(), "u-switch", "org1")
	if err != nil || len(perms1) == 0 {
		t.Fatalf("expected perms for org1")
	}
	perms2, err := svc.GetUserPermissions(context.Background(), "u-switch", "org2")
	if err != nil || len(perms2) == 0 {
		t.Fatalf("expected perms for org2")
	}
	if perms1[0] == perms2[0] {
		t.Fatalf("expected different permissions after org switch")
	}
}

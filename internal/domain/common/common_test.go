package common

import (
	"testing"
)

func TestIDs(t *testing.T) {
	var pid ProjectID = "proj-1"
	var wid WorkspaceID = "work-1"

	// This is primarily a compile-time check that types are distinct.
	// If they were the same, we could assign them. Go enforces strict types here.
	if string(pid) == string(wid) {
		t.Fatalf("ids should not match")
	}
}

func TestEnums(t *testing.T) {
	if EnvironmentTypeProduction != "production" {
		t.Fatalf("expected production")
	}
	if ExecutorTypeDocker != "docker" {
		t.Fatalf("expected docker")
	}
}

func TestValidation(t *testing.T) {
	if err := ValidateName(""); err == nil {
		t.Fatalf("expected error for empty name")
	}
	if err := ValidatePath(" "); err == nil {
		t.Fatalf("expected error for blank path")
	}
	if err := ValidateEndpoint(""); err == nil {
		t.Fatalf("expected error for empty endpoint")
	}
	if err := ValidateVersion("   "); err == nil {
		t.Fatalf("expected error for blank version")
	}
	if err := ValidateReference(""); err == nil {
		t.Fatalf("expected error for empty reference")
	}
}

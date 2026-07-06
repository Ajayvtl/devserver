package valueobjects

import (
	"testing"
)

func TestName(t *testing.T) {
	_, err := NewName("ValidName")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = NewName("   ")
	if err == nil {
		t.Fatalf("expected error for empty name, got nil")
	}
}

func TestPath(t *testing.T) {
	_, err := NewPath("/opt/devserver")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = NewPath("")
	if err == nil {
		t.Fatalf("expected error for empty path, got nil")
	}
}

func TestEndpoint(t *testing.T) {
	_, err := NewEndpoint("localhost:3000")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = NewEndpoint("")
	if err == nil {
		t.Fatalf("expected error for empty endpoint, got nil")
	}
}

func TestVersion(t *testing.T) {
	_, err := NewVersion("1.0.0")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = NewVersion(" ")
	if err == nil {
		t.Fatalf("expected error for empty version, got nil")
	}
}

func TestReference(t *testing.T) {
	_, err := NewReference("github.com/org/repo")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = NewReference("")
	if err == nil {
		t.Fatalf("expected error for empty reference, got nil")
	}
}

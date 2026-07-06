package contracts

import "github.com/Ajayvtl/devserver/internal/domain/common"

// ApplicationResult represents the canonical output of any Use Case or Application Service.
// It wraps the raw domain response with additional application-level metadata.
type ApplicationResult struct {
	EntityID common.ProjectID // Or any generalized ID if needed
	Success  bool
	Message  string
	Warnings []string
}

package valueobjects

import "github.com/Ajayvtl/devserver/internal/domain/common"

// Name is a validated value object.
type Name struct {
	value string
}

// NewName creates a new validated Name value object.
func NewName(v string) (Name, error) {
	if err := common.ValidateName(v); err != nil {
		return Name{}, err
	}
	return Name{value: v}, nil
}

// String returns the underlying string value.
func (n Name) String() string {
	return n.value
}

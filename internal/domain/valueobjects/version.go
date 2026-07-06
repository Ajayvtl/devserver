package valueobjects

import "github.com/Ajayvtl/devserver/internal/domain/common"

type Version struct {
	value string
}

func NewVersion(v string) (Version, error) {
	if err := common.ValidateVersion(v); err != nil {
		return Version{}, err
	}
	return Version{value: v}, nil
}

func (v Version) String() string {
	return v.value
}

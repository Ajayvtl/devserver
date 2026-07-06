package valueobjects

import "github.com/Ajayvtl/devserver/internal/domain/common"

type Path struct {
	value string
}

func NewPath(v string) (Path, error) {
	if err := common.ValidatePath(v); err != nil {
		return Path{}, err
	}
	return Path{value: v}, nil
}

func (p Path) String() string {
	return p.value
}

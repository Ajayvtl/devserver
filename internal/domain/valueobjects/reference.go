package valueobjects

import "github.com/Ajayvtl/devserver/internal/domain/common"

type Reference struct {
	value string
}

func NewReference(v string) (Reference, error) {
	if err := common.ValidateReference(v); err != nil {
		return Reference{}, err
	}
	return Reference{value: v}, nil
}

func (r Reference) String() string {
	return r.value
}

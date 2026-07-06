package valueobjects

import "github.com/Ajayvtl/devserver/internal/domain/common"

type Endpoint struct {
	value string
}

func NewEndpoint(v string) (Endpoint, error) {
	if err := common.ValidateEndpoint(v); err != nil {
		return Endpoint{}, err
	}
	return Endpoint{value: v}, nil
}

func (e Endpoint) String() string {
	return e.value
}

package common

import "errors"

var (
	ErrOptimisticLock = errors.New("optimistic lock failed: resource modified")
)

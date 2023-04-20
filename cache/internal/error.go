package internal

import (
	"errors"
	"fmt"
)

var (
	ErrKeyNotFound = errors.New("cache: 未找到key")
)

func NewErrKeyNotFound(key string) error {
	return fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

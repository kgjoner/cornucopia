package cache

import (
	"context"
	"fmt"
	"time"
)

// Error when cache key does not exist. It does not apply if key does exist but it refers to a zero value.
var ErrNil = fmt.Errorf("cache: nil")

type Pool interface {
	NewDAO(context.Context) DAO
	Close() error
}

type DAO interface {
	CacheJson(key string, v interface{}, duration time.Duration) error
	GetJson(key string, v interface{}) error
	Clear(key string)
}

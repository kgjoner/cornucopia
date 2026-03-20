package cache

import (
	"context"
	"fmt"
	"time"
)

// Error when cache key does not exist. It does not apply if key does exist but it refers to a zero value.
var ErrNil = fmt.Errorf("cache: nil")

type Pool interface {
	NewStore(context.Context) Store
	Close() error
}

type Store interface {
	CacheJSON(key string, v interface{}, duration time.Duration) error
	GetJSON(key string, v interface{}) error
	Clear(key string)
}

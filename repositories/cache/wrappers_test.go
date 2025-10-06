package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgjoner/cornucopia/v2/repositories/cache"
	"github.com/kgjoner/cornucopia/v2/repositories/cache/memorydb"
	"github.com/stretchr/testify/assert"
)

var cachePool, _ = memorydb.NewPool()
var cacheRepo = cachePool.NewDAO(context.Background())
var sleepTime = 1 * time.Second

func TestRunWithCache(t *testing.T) {
	start := time.Now()

	var mockingA struct{}
	res, err := cache.RunWithCache[result](cacheRepo, 5*time.Minute, NewResult)(1, "foo", &mockingA)
	assert.Nil(t, err)
	assert.Equal(t, res.ID, 1)
	assert.GreaterOrEqual(t, time.Since(start), sleepTime)

	start = time.Now()

	var mockingB struct{}
	res, err = cache.RunWithCache[result](cacheRepo, 5*time.Minute, NewResult)(1, "foo", &mockingB)
	assert.Nil(t, err)
	assert.Equal(t, res.ID, 1)
	assert.Equal(t, res.Name, "foo")
	assert.Less(t, time.Since(start), sleepTime)
}

type result struct {
	ID   int
	Name string
}

func NewResult(id int, name string, data interface{}) (*result, error) {
	time.Sleep(sleepTime)
	return &result{
		id,
		name,
	}, nil
}

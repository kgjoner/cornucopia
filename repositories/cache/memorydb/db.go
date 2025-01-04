package memorydb

import (
	"context"

	"github.com/kgjoner/cornucopia/repositories/cache"
)

// Simple in memory cache for tests or pocs. It does not implement duration.
type Pool struct {
	data map[string]string
}

func NewPool() (*Pool, error) {
	return &Pool{
		data: map[string]string{},
	}, nil
}

type DAO struct {
	ctx  context.Context
	data map[string]string
}

func (p Pool) NewDAO(ctx context.Context) cache.DAO {
	return &DAO{
		ctx:  ctx,
		data: p.data,
	}
}

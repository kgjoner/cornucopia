package memorydb

import (
	"context"

	"github.com/kgjoner/cornucopia/v2/cache"
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

func (p *Pool) Close() error {
	p.data = nil
	return nil
}

type Store struct {
	ctx  context.Context
	data map[string]string
}

func (p *Pool) NewDAO(ctx context.Context) cache.Store {
	return &Store{
		ctx:  ctx,
		data: p.data,
	}
}

package memorydb

import (
	"context"
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

func (p Pool) NewDAO(ctx context.Context) *DAO {
	return &DAO{
		ctx:  ctx,
		data: p.data,
	}
}

package redisdb

import (
	"context"
	"fmt"

	"github.com/kgjoner/cornucopia/v2/repositories/cache"
	"github.com/redis/go-redis/v9"
)

type Pool struct {
	url string
	db  *redis.Client
}

func NewPool(url string) (*Pool, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("cacherepo: unable to parse redis url: %v", err)
	}

	return &Pool{
		url,
		redis.NewClient(opt),
	}, nil
}

func (p *Pool) Close() error {
	if p.db == nil {
		return nil
	}
	if err := p.db.Close(); err != nil {
		return fmt.Errorf("cacherepo: unable to close redis connection: %v", err)
	}

	p.db = nil
	return nil
}

func (p *Pool) DatabaseURL() string {
	return p.url
}

func (p *Pool) Client() *redis.Client {
	return p.db
}

type DAO struct {
	ctx context.Context
	db  *redis.Client
}

func (p Pool) NewDAO(ctx context.Context) cache.DAO {
	return &DAO{
		ctx: ctx,
		db:  p.db,
	}
}

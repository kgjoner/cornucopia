package cacherepo

import (
	"context"
	"fmt"

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

func (p Pool) DatabaseUrl() string {
	return p.url
}

func (p Pool) Client() *redis.Client {
	return p.db
}

type Queries struct {
	ctx context.Context
	db  *redis.Client
}

func (p Pool) NewQueries(ctx context.Context) *Queries {
	return &Queries{
		ctx: ctx,
		db:  p.db,
	}
}


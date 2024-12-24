package redisdb

import (
	"encoding/json"
	"time"

	"github.com/kgjoner/cornucopia/repositories/cache"
	"github.com/redis/go-redis/v9"
)

func (q DAO) CacheJson(key string, v interface{}, duration time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return q.db.Set(q.ctx, key, string(data), duration).Err()
}

func (q DAO) GetJson(key string, v interface{}) error {
	jsonData, err := q.db.Get(q.ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	} else if err == redis.Nil {
		return cache.ErrNil
	}

	err = json.Unmarshal([]byte(jsonData), v)
	if err != nil {
		return err
	}

	return err
}

func (q DAO) Clear(key string) {
	q.db.Del(q.ctx, key)
}
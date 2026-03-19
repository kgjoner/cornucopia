package memorydb

import (
	"encoding/json"
	"time"

	"github.com/kgjoner/cornucopia/v2/cache"
)

func (q Store) CacheJSON(key string, v interface{}, duration time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	q.data[key] = string(data)
	return nil
}

func (q Store) GetJSON(key string, v interface{}) error {
	jsonData, exists := q.data[key]
	if !exists {
		return cache.ErrNil
	}

	err := json.Unmarshal([]byte(jsonData), v)
	if err != nil {
		return err
	}

	return err
}

func (q Store) Clear(key string) {
	delete(q.data, key)
}

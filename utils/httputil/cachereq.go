package httputil

import (
	"net/http"
	"time"

	cacherepo "github.com/kgjoner/cornucopia/repositories/cache"
)

func GetWithCache(url string, data any, cacheRepo cacherepo.Queries, cacheDuration time.Duration) error {
	key := "get:" + url
	err := cacheRepo.GetJson(key, data)
	if err != cacherepo.ErrNil {
		// hit the cache (err == nil) or got a cache internal error
		return err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	_, err = DoReq(client, req, data)
	if err != nil {
		return err
	}

	err = cacheRepo.CacheJson(key, data, cacheDuration)
	if err != nil {
	 return err
	}

	return nil
}

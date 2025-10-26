package cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

// Cache is the interface for the cache
// get retrieve a value from the cache
// set store a value in the cache
// delete remove a value from the cache
type ICache interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
	Ping(ctx context.Context) error
}

// GetOrLoad utility function to retrieve a value from the cache, if not found, load it using the loader function
// if the cache is nil, it will load the value using the loader function
func GetOrLoad[T any](ctx context.Context, cache ICache, key string, expiration time.Duration, loadFunc func() (T, error)) (T, error) {
	var value T
	var zero T
	var err error

	// if cache is nil or key is empty, return the result of loadFunc
	if cache == nil || key == "" {
		return loadFunc()
	}

	err = cache.Get(ctx, key, &value)
	if err == nil {
		return value, nil
	}

	// prevent multiple requests update the cache for the same key
	result, err, _ := group.Do(key, func() (interface{}, error) {
		// security validate if not is cached
		var cached T
		if e := cache.Get(ctx, key, &cached); e == nil {
			return cached, nil
		}

		// if cache miss, load the value using the loader function
		_value, _err := loadFunc()
		if _err != nil {
			return _value, _err
		}
		cache.Set(ctx, key, _value, expiration)
		return _value, nil
	})

	if err != nil {
		return zero, err
	}

	return result.(T), nil
}

package cache

import (
	"context"
	"time"
)

// Cache is the interface for the cache
// get retrieve a value from the cache
// set store a value in the cache
// delete remove a value from the cache
type ICache interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

// GetOrLoad utility function to retrieve a value from the cache, if not found, load it using the loader function
// if the cache is nil, it will load the value using the loader function
func GetOrLoad[T any](ctx context.Context, cache ICache, key string, expiration time.Duration, loadFunc func() (T, error)) (T, error) {
	var value T

	// if cache is nil, return the result of loadFunc
	if cache == nil {
		return loadFunc()
	}

	err := cache.Get(ctx, key, &value)
	if err == nil {
		return value, nil
	}

	// if cache miss, load the value using the loader function
	value, err = loadFunc()
	if err != nil {
		return value, err
	}

	// set the value in the cache
	cache.Set(ctx, key, value, expiration)
	return value, nil
}

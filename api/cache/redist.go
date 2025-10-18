package cache

import (
	"api/config"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Reddis struct {
	client *redis.Client
}

func NewReddis() (*Reddis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Cache().Addr,
		Password:     config.Cache().Password,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 1,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Reddis{
		client: client,
	}, nil
}

// Get retrieves a value from the cache
// If the key does not exist, returns nil
func (r *Reddis) Get(ctx context.Context, key string, value interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

// Set sets a value in the cache
func (r *Reddis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Delete deletes a key from the cache
func (r *Reddis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Reddis) Close() error {
	return r.client.Close()
}

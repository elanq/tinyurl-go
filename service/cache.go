package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elanq/pastebin-go/connection"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(context.Context, string) (interface{}, error)
	Evict(ctx context.Context, key string) error
}

type redisCache struct {
	client *connection.Redis
}

// Evict implements Cache.
func (r *redisCache) Evict(ctx context.Context, key string) error {
	res := r.client.Client().Del(ctx, key)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

// Get implements Cache.
func (r *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	res := r.client.Client().Get(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.Val(), nil
}

// Set implements Cache.
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	res := r.client.Client().Set(ctx, key, v, ttl)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func NewRedisCache(r *connection.Redis) Cache {
	return &redisCache{
		client: r,
	}
}

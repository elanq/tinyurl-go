package connection

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	red *redis.Client
}

func (r *Redis) Client() *redis.Client {
	return r.red
}

func NewRedis() *Redis {
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	r := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	err := r.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	return &Redis{
		red: r,
	}
}

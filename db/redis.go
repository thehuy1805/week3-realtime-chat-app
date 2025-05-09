package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
    Client *redis.Client
}

func NewRedisDB(addr string) (*RedisDB, error) {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        return nil, err
    }
    return &RedisDB{Client: client}, nil
}

func (r *RedisDB) Close() {
    r.Client.Close()
}
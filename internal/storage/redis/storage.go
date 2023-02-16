package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewDB(addr, password string, db int) (conn *redis.Client, err error) {
	conn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err = conn.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}
	return conn, nil
}

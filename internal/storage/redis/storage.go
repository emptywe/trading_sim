package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string
	Port     string
	Password string
	Db       string
}

func NewDB(cfg Config) (conn *redis.Client, err error) {
	db, err := strconv.ParseInt(cfg.Db, 10, 8)
	if err != nil {
		return nil, err
	}
	conn = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       int(db),
	})
	if _, err = conn.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}
	return conn, nil
}

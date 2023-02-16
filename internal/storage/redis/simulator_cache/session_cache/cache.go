package session_cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Session interface {
	Create(ctx context.Context, username, ssid, token string) error
	Read(ctx context.Context, username, ssid, token string) error
	Update(ctx context.Context, username, ssid, token string) error
	Delete(ctx context.Context, username string) error
}

type Cache struct {
	Session
}

func NewCache(db *redis.Client) *Cache {
	return &Cache{
		Session: NewSessionCache(db),
	}
}

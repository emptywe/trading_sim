package session_cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/emptywe/trading_sim/pkg/session"
)

type SessionCache struct {
	db *redis.Client
}

func NewSessionCache(db *redis.Client) *SessionCache {
	return &SessionCache{db: db}
}

func (c *SessionCache) Create(ctx context.Context, username, ssid, token string) error {
	key := fmt.Sprintf("session_cache%s_%s", username, ssid)
	err := c.db.Set(ctx, key, token, session.ExpireSession).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *SessionCache) Read(ctx context.Context, username, ssid, token string) error {
	result, err := c.db.Get(ctx, fmt.Sprintf("session_cache%s_%s", username, ssid)).Result()
	if err != nil {
		return err
	}

	if result != token {
		return errors.New("invalid session_cache or token")
	}

	return nil
}

func (c *SessionCache) Update(ctx context.Context, username, ssid, token string) error {
	key := fmt.Sprintf("session_cache%s_%s", username, ssid)

	err := c.db.Set(ctx, key, token, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *SessionCache) Delete(ctx context.Context, username string) error {
	var keys []string
	var cursor uint64

	for {
		var err error
		var key []string
		key, cursor, err = c.db.Scan(ctx, cursor, fmt.Sprintf("session_cache%s_*", username), 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, key...)
		if cursor == 0 {
			break
		}
	}
	for _, k := range keys {
		_, err := c.db.Del(ctx, k).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

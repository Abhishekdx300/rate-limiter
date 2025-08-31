package limiter

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed token_bucket.lua
var luaScript string

type RateLimiter struct {
	client *redis.Client
	script *redis.Script
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
		script: redis.NewScript(luaScript),
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string, limit int, rate float64) (bool, error) {
	args := []any{limit, rate, float64(time.Now().UnixNano()) / 1e9}

	res, err := rl.script.Run(ctx, rl.client, []string{key}, args...).Result()
	if err != nil {
		return false, err
	}

	resultSlice := res.([]any)

	allowed := (resultSlice[0].(int64) == 1)

	return allowed, nil
}

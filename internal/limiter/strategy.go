package limiter

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"time"
)

type Strategy string

const (
	TokenBucket          Strategy = "token-bucket"
	LeakingBucket        Strategy = "leaking-bucket"
	FixedWindow          Strategy = "fixed-window"
	SlidingWindowLog     Strategy = "sliding-window-log"
	SlidingWindowCounter Strategy = "sliding-window-counter"
)

type Limiter interface {
	AllowRequest(ctx context.Context, key string) (bool, int, error)
}

type FixedWindowLimiter struct {
	Redis  rueidis.Client
	Limit  int           // Maximum request per window
	Window time.Duration // e.g. 1 minute
	Prefix string        // Redis key prefix
}

// AllowRequest implements fixed window rate limiting with atomic Redis logic
func (l *FixedWindowLimiter) AllowRequest(ctx context.Context, key string) (bool, int, error) {
	fullKey := fmt.Sprintf("%s%s", l.Prefix, key)

	// Lua script for atomic increment and expiration
	script := `
	local current
	current = redis.call("INCR", KEYS[1])
	if tonumber(current) == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end
	return current
	`
	ttlSeconds := int(l.Window.Seconds())
	cmd := l.Redis.B().
		Eval().
		Script(script).
		Numkeys(1).
		Key(fullKey).
		Arg(fmt.Sprintf("%d", ttlSeconds)).
		Build()
	result := l.Redis.Do(ctx, cmd)
	currentCount, err := result.AsInt64()
	if err != nil {
		return false, 0, err
	}
	if currentCount > int64(l.Limit) {
		// Estimate retry-after using window (not 100% accurate without TTL read)
		return false, ttlSeconds, nil
	}
	return true, ttlSeconds, nil
}

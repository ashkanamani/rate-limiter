package limiter

import (
	"github.com/redis/rueidis"
	"time"
)

func NewLimiter(redis rueidis.Client, limit int, window time.Duration, prefix string, strategy Strategy) Limiter {
	switch strategy {
	case FixedWindow:
		return &FixedWindowLimiter{
			Redis:  redis,
			Limit:  limit,
			Window: window,
			Prefix: prefix,
		}
	default:
		panic("unsupported rate limiter strategy: " + string(strategy))
	}
}

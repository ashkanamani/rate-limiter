package limiter

import (
	"context"
	"github.com/ashkanamani/rate-limiter/config"
	"github.com/redis/rueidis"
	"time"
)

func InitRedis(cfg *config.Config) (rueidis.Client, error) {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{cfg.RedisAddr},
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return nil, err
	}
	return client, err
}

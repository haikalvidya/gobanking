package app

import (
	"context"
	"gobanking/pkg/redis"
	"time"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
)

func (a *app) connectRedis(ctx context.Context) error {

	retryOptions := []retry.Option{
		retry.Attempts(a.cfg.Timeouts.RedisInitRetryCount),
		retry.Delay(time.Duration(a.cfg.Timeouts.RedisInitMilliseconds) * time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			a.log.Errorf("retry connect redis err: %v", err)
		}),
	}

	return retry.Do(func() error {
		redisConn, err := redis.NewRedisConn(&ctx, a.cfg.Redis)
		if err != nil {
			return errors.Wrap(err, "redis.NewRedisConn")
		}
		a.redisConn = redisConn
		return nil
	}, retryOptions...)
}

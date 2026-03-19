package redis

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	_redis "github.com/redis/go-redis/v9"
)

func Lock[T any](ctx context.Context, key string, f func() T) (t *T, err error) {
	u := uuid.New().String()
	res, err := RDB.SetArgs(ctx, key, u, _redis.SetArgs{
		Mode: "NX",
		TTL:  5 * time.Second,
	}).Result()
	if err != nil {
		log.Printf("lock not acquired, key: %#v, error: %#v\n", key, err)
		return nil, err
	}
	if res != "OK" {
		log.Printf("lock not acquired, key: %#v\n", key)
		return nil, errors.New("lock not acquired")
	}
	defer func() {
		_, err := _redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`).Run(ctx, RDB, []string{key}, u).Result()
		//res, err = RDB.Get(ctx, key).Result()
		if err != nil {
			log.Printf("unlock failed, key=%s err=%v", key, err)
			//return
		}
		//if res == u {
		//	RDB.Del(ctx, key)
		//}
	}()
	v := f()
	return &v, nil
}

func Acquire(ctx context.Context, key string, timeout time.Duration, acquireTimeout time.Duration) (string, error) {
	u := uuid.New().String()
	attempts := 0
	deadline := time.Now().Add(acquireTimeout)
	for time.Now().Before(deadline) {
		res, err := RDB.SetArgs(ctx, key, u, _redis.SetArgs{
			Mode: "NX",
			TTL:  timeout,
		}).Result()
		if err != nil {
			log.Printf("lock not acquired, key: %#v, error: %#v\n", key, err)
			return "", err
		}
		if res == "OK" {
			return u, nil
		}
		time.Sleep(calculateExponentialBackoff(attempts))
		attempts++
	}
	return "", errors.New("当前正忙，请稍后再试")
}

func calculateExponentialBackoff(attempts int) time.Duration {
	backoff := 100 * (1 << attempts)
	if backoff > 1000 {
		backoff = 1000
	}
	return time.Duration(backoff) * time.Millisecond
}

func Release(ctx context.Context, key string, identifier string) {
	_, err := _redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`).Run(ctx, RDB, []string{key}, identifier).Result()
	if err != nil {
		log.Printf("release failed, key=%s identifier=%v err=%v", key, identifier, err)
		//return
	}
}

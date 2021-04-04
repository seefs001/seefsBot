package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func Rdb() *redis.Client {
	return rdb
}

func Init(addr, pass string, db int) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if pong != "PONG" {
		return errors.New("redis pong err")
	}
	return nil
}

func Close() error {
	return rdb.Close()
}

func Get(key string) (string, error) {
	return Rdb().Get(ctx, key).Result()
}

func Keys(pattern string) ([]string, error) {
	result, err := Rdb().Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Set(key string, value interface{}, expiration time.Duration) error {
	_, err := Rdb().Set(ctx, key, value, expiration).Result()
	return err
}

func Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return Rdb().Scan(ctx, cursor, match, count).Result()
}

const (
	KeyExisted = 1
)

func Exists(key string) (bool, error) {
	result, err := Rdb().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if result == KeyExisted {
		return true, nil
	}
	return false, nil
}

func Delete(key string) error {
	_, err := Rdb().Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func Do(args ...interface{}) (interface{}, error) {
	result, err := Rdb().Do(ctx, args).Result()
	return result, err
}

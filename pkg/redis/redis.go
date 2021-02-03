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

func checkRdb() {
	panic("rdb is not initialized")
}

func Rdb() *redis.Client {
	checkRdb()
	return rdb
}

// Init initialize redis client
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

// Get get value from redis
func Get(key string) (string, error) {
	checkRdb()
	return rdb.Get(ctx, key).Result()
}

// Keys ...
func Keys(pattern string) ([]string, error) {
	checkRdb()
	result, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Set set value to redis
func Set(key string, value interface{}, expiration time.Duration) error {
	checkRdb()
	_, err := rdb.Set(ctx, key, value, expiration).Result()
	return err
}

// Scan ...
func Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	checkRdb()
	return rdb.Scan(ctx, cursor, match, count).Result()
}

const (
	KeyExisted = 1
)

// Exists exist in redis
func Exists(key string) (bool, error) {
	checkRdb()
	result, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if result == KeyExisted {
		return true, nil
	}
	return false, nil
}

// Delete delete from redis
func Delete(key string) error {
	checkRdb()
	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Do ...
func Do(args ...interface{}) (interface{}, error) {
	checkRdb()
	result, err := rdb.Do(ctx, args).Result()
	return result, err
}

// Publish publish to channel
func Publish(channel string, message interface{}) error {
	checkRdb()
	return rdb.Publish(ctx, channel, message).Err()
}

// Subscribe subscribe a channel
func Subscribe(channel string) *redis.PubSub {
	checkRdb()
	return rdb.Subscribe(ctx, channel)
}

// SubscribeChan subscribe a channel,return <- chan
func SubscribeChan(channel string) <-chan *redis.Message {
	checkRdb()
	return rdb.Subscribe(ctx, channel).Channel()
}

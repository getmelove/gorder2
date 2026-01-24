package redis

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

// CacheClient 定义了高级缓存操作的接口
type CacheClient interface {
	// GetOrSet 获取缓存，如果不存在则执行 fetch 函数获取并写入缓存。
	// 包含防击穿（singleflight）、防穿透（空值缓存）、防雪崩（随机 TTL）功能。
	GetOrSet(ctx context.Context, key string, ttl time.Duration, fetch func(ctx context.Context) (string, error)) (string, error)
}

type cacheClient struct {
	client *redis.Client
	sf     singleflight.Group
}

var (
	// ErrCacheMiss 表示缓存未命中
	ErrCacheMiss = errors.New("cache miss")
	// emptyValue 是用于防止缓存穿透的空值占位符
	emptyValue = "*"
)

// NewCacheClient 创建一个新的 CacheClient 实例
func NewCacheClient(client *redis.Client) CacheClient {
	return &cacheClient{
		client: client,
	}
}

// DefaultCacheClient 返回使用默认 LocalClient 的 CacheClient
func DefaultCacheClient() CacheClient {
	return NewCacheClient(LocalClient())
}

func (c *cacheClient) GetOrSet(ctx context.Context, key string, ttl time.Duration, fetch func(ctx context.Context) (string, error)) (string, error) {
	// 1. 尝试从缓存获取
	val, err := c.get(ctx, key)
	if err == nil {
		if val == emptyValue {
			return "", nil // 命中空值缓存，直接返回空字符串，防止穿透
		}
		return val, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		logrus.Warnf("redis get error: %v", err)
		// 如果 Redis 报错（非 Miss），我们通常应该降级去查 DB，而不是直接报错
	}

	// 2. 缓存未命中，使用 Singleflight 防止击穿
	// Do 方法确保对同一个 key 的并发调用，只有一个会真正执行 fn
	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		// 再次检查缓存（Double Check），防止在上锁期间已被其他请求写入
		if val, err := c.get(ctx, key); err == nil {
			return val, nil
		}

		// 执行真正的获取逻辑（查询 DB 或 RPC）
		fetchedVal, fetchErr := fetch(ctx)
		if fetchErr != nil {
			// 如果获取失败，不缓存错误（或者可以根据业务决定是否缓存错误）
			return "", fetchErr
		}

		// 3. 处理防穿透：如果结果为空，也进行缓存（时间设置短一点）
		storeVal := fetchedVal
		realTTL := ttl
		if fetchedVal == "" {
			storeVal = emptyValue
			realTTL = 5 * time.Minute // 空值缓存时间短一些
		} else {
			// 4. 处理防雪崩：增加随机 TTL 抖动 (0-10%)
			// 只有正常值才需要抖动，空值通常短且固定
			jitter := time.Duration(rand.Int63n(int64(ttl / 10)))
			realTTL += jitter
		}

		// 写入缓存
		if err := c.set(ctx, key, storeVal, realTTL); err != nil {
			logrus.Warnf("redis set error: %v", err)
		}

		return storeVal, nil
	})

	if err != nil {
		return "", err
	}

	result := v.(string)
	if result == emptyValue {
		return "", nil
	}
	return result, nil
}

// get 是内部方法，简单封装了 Redis Get
func (c *cacheClient) get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCacheMiss
	}
	return val, err
}

// set 是内部方法，简单封装了 Redis Set
func (c *cacheClient) set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

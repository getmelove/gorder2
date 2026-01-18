package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// SetNX 封装了 Redis 的 setnx (set if not exists) 命令。
// 如果 key 不存在，则设置 key 的值为 value，并设置过期时间 ttl。
// 返回 err 为 nil 表示操作成功。
//
// 参数:
//   - ctx: 包含上下文信息，用于控制请求超时和传递值。
//   - client: Redis 客户端实例。
//   - key: 要设置的键。
//   - value: 要设置的值。
//   - ttl: 键的过期时间 (Time To Live)。
func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl time.Duration) (err error) {
	now := time.Now()
	// 使用 defer 来记录日志，无论函数是正常返回还是发生错误都会执行。
	// 这是一种常见的"切面编程"手法，用于统一处理监控和日志。
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start": now,
			"key":   key,
			"value": value,
			"err":   err,                            // 记录操作产生的错误（如果有）
			"cost":  time.Since(now).Milliseconds(), // 记录操作耗时，单位毫秒
		})
		if err == nil {
			l.Info("redis_setnx_success")
		} else {
			l.Warn("redis_setnx_error")
		}
	}()

	if client == nil {
		return errors.New("redis client is nil")
	}
	// 调用 go-redis 库的 SetNX 方法
	_, err = client.SetNX(ctx, key, value, ttl).Result()
	return err
}

// Del 封装了 Redis 的 del 命令，用于删除指定的 key。
//
// 参数:
//   - ctx: 上下文。
//   - client: Redis 客户端实例。
//   - key: 要删除的键。
func Del(ctx context.Context, client *redis.Client, key string) (err error) {
	now := time.Now()
	// 同样使用 defer 记录操作结束后的日志和耗时。
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start": now,
			"key":   key,
			"err":   err,
			"cost":  time.Since(now).Milliseconds(),
		})
		if err == nil {
			l.Info("redis_del_success")
		} else {
			l.Warn("redis_del_error")
		}
	}()

	if client == nil {
		return errors.New("redis client is nil")
	}
	// 调用 go-redis 库的 Del 方法
	_, err = client.Del(ctx, key).Result()
	return err
}

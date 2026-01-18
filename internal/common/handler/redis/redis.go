package redis

import (
	"fmt"
	"time"

	"github.com/getmelove/gorder2/internal/common/handler/factory"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	// confName 定义了配置文件中 Redis 配置的基本路径名称（例如在 yaml 文件中的 "redis" 块）
	confName = "redis"
	// localSupplier 是默认的本地 Redis 连接标识符，对应配置文件中 redis.local
	localSupplier = "local"
)

var (
	// singleton 是一个单例工厂实例，用于管理 Redis 客户端的创建和获取。
	// 它使用了 factory.NewSingleton 并传入 supplier 函数作为创建者。
	// 这样可以确保对于同一个配置名的 Redis 客户端，只会创建一个实例。
	singleton = factory.NewSingleton(supplier)
)

// Init 初始化 Redis 模块。
// 它会从配置文件中读取 "redis" 下的所有配置项，并预先初始化它们。
// 这里的 conf 是一个 map，key 是配置名称（如 local, cache 等），value 是具体的配置内容。
func Init() {
	conf := viper.GetStringMap(confName)
	for supplyName := range conf {
		Client(supplyName)
	}
}

// LocalClient 获取名为 "local" 的 Redis 客户端实例。
// 这是一个便捷方法，通常用于获取默认的 Redis 连接。
func LocalClient() *redis.Client {
	return Client(localSupplier)
}

// Client 根据提供的名称获取对应的 Redis 客户端实例。
// 它通过单例工厂获取实例，如果实例不存在，工厂会自动调用 supplier 创建。
// 返回值被断言为 *redis.Client 类型。
func Client(name string) *redis.Client {
	return singleton.Get(name).(*redis.Client)
}

// supplier 是 Singleton 工厂的回调函数，用于实际创建 Redis 客户端。
// key 参数是配置名称（例如 "local"）。
func supplier(key string) any {
	// 拼接完整的配置路径，例如 "redis.local"
	confKey := confName + "." + key

	// Section 定义了 Redis 配置的结构体，映射配置文件中的字段。
	type Section struct {
		IP           string        `mapstructure:"ip"`            // Redis服务器IP
		Port         string        `mapstructure:"port"`          // Redis服务器端口
		PoolSize     int           `mapstructure:"pool_size"`     // 连接池大小
		MaxConn      int           `mapstructure:"max_conn"`      // 最大连接数
		ConnTimeout  time.Duration `mapstructure:"conn_timeout"`  // 连接超时时间
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 读取超时时间
		WriteTimeout time.Duration `mapstructure:"write_timeout"` // 写入超时时间
	}
	var c Section
	// 使用 Viper 解析配置到结构体中
	if err := viper.UnmarshalKey(confKey, &c); err != nil {
		panic(err) // 如果配置解析失败，直接 panic，因为这是启动时的关键错误
	}

	// 使用 go-redis 库创建并返回一个新的 Redis 客户端
	return redis.NewClient(&redis.Options{
		Network:         "tcp",
		Addr:            fmt.Sprintf("%s:%s", c.IP, c.Port),
		PoolSize:        c.PoolSize,                       // 设置连接池大小
		MaxActiveConns:  c.MaxConn,                        // 设置最大活跃连接数
		ConnMaxLifetime: c.ConnTimeout * time.Millisecond, // 将配置的数值转换为时长（毫秒）
		ReadTimeout:     c.ReadTimeout * time.Millisecond,
		WriteTimeout:    c.WriteTimeout * time.Millisecond,
	})
}

package red

import (
	"commmunity/app/config"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

type Redis struct {
	redis   *redis.Client
	context context.Context
}

func NewRedis(redis *redis.Client, ctx context.Context) *Redis {
	return &Redis{
		redis:   redis,
		context: ctx,
	}
}
func ConnectRedis() (*redis.Client, context.Context) {
	var red config.RedisConfig
	red.Host = viper.GetString("redis.host")
	red.Port = viper.GetString("redis.port")
	addr := fmt.Sprintf("%s:%s", red.Host, red.Port)
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx := context.Background()
	return client, ctx
}

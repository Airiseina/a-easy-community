package config

import (
	"commmunity/app/zlog"

	"github.com/spf13/viper"
)

type MysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
type RedisConfig struct {
	Host string
	Port string
}

func GetConfig() {
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", "3306")
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "123456")
	viper.SetDefault("mysql.name", "user")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("jwtKey", "EL PSY KONGROO")
	viper.SetDefault("jwtRefreshKey", "Steins Gate")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		zlog.Info("未找到配置文件，使用默认值")
	}
}

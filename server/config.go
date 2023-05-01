package server

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	RedisPort            int           `mapstructure:"REDIS_PORT"`
	RedisQueue           bool          `mapstructure:"REDIS_QUEUE"`
	TokenHashKey         string        `mapstructure:"TOKEN_HASH_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDueation time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	RootUsername         string        `mapstructure:"ROOT_LOGIN_USERNAME"`
	RootPwd              string        `mapstructure:"ROOT_LOGIN_PWD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}

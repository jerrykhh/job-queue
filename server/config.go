package server

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	RedisPort            int           `mapstructure:"REDIS_PORT"`
	TokenHashKey         string        `mapstructure:"TOKEN_HASH_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDueation time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}

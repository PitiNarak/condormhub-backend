package redis

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Uri                   string `env:"URI"`
	AccessTokenExpireHrs  int    `env:"ACCESS_TOKEN_EXPIRE_HRS"`
	RefreshTokenExpireHrs int    `env:"REFRESH_TOKEN_EXPIRE_HRS"`
}

type Redis struct {
	client *redis.Client
}

func New(config Config) (*Redis, error) {
	options, err := redis.ParseURL(config.Uri)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	return &Redis{
		client: client,
	}, nil
}

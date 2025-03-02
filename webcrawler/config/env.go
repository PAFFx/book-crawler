package config

import (
	"log"

	"github.com/Netflix/go-env"
)

type EnvVariables struct {
	CrawlerThreads int `env:"CRAWLER_THREADS,default=8"`

	RedisHost     string `env:"REDIS_HOST,default=localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD,default=1q2w3e4r"`
	RedisDB       int    `env:"REDIS_DB,default=0"`
}

func GetEnv() (*EnvVariables, error) {
	var cfg EnvVariables
	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		log.Panicf("Error unmarshalling environment variables: %v", err)
		return nil, err
	}
	return &cfg, nil
}

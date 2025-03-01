package redis

import (
	"github.com/gocolly/redisstorage"

	"book-search/webcrawler/config"
)

var storage *redisstorage.Storage

func InitStorage() error {
	env, err := config.GetEnv()
	if err != nil {
		return err
	}
	storage = &redisstorage.Storage{
		Address:  env.RedisHost,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
		Prefix:   "webcrawler:",
	}

	return nil
}

func GetStorage() *redisstorage.Storage {
	return storage
}

func CloseStorageClient() error {
	return storage.Client.Close()
}

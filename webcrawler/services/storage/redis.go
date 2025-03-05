package storage

import (
	"github.com/gocolly/redisstorage"

	"book-search/webcrawler/config"
)

func GetStorage(prefix string) (*redisstorage.Storage, error) {
	env, err := config.GetEnv()
	if err != nil {
		return nil, err
	}

	instance := &redisstorage.Storage{
		Address:  env.RedisHost,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
		Prefix:   prefix,
	}

	return instance, nil
}

func CloseStorageClient(instance *redisstorage.Storage) error {
	return instance.Client.Close()
}

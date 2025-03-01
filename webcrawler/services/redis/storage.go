package redis

import (
	"log"
	"sync"

	"github.com/gocolly/redisstorage"

	"book-search/webcrawler/config"
)

var once sync.Once
var instance *redisstorage.Storage

func GetStorage() (*redisstorage.Storage, error) {
	env, err := config.GetEnv()
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		log.Println("Initializing Redis storage backend")
		instance = &redisstorage.Storage{
			Address:  env.RedisHost,
			Password: env.RedisPassword,
			DB:       env.RedisDB,
			Prefix:   "webcrawler:",
		}
	})

	return instance, nil
}

func CloseStorageClient() error {
	return instance.Client.Close()
}

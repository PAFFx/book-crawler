package main

import (
	"context"
	"log"

	"book-search/webcrawler/crawler"
	"book-search/webcrawler/services/redis"
)

func main() {
	ctx := context.Background()

	// Init services
	if err := redis.InitRedisClient(ctx); err != nil {
		log.Fatal(err)
	}
	defer redis.CloseRedisClient()
	log.Println("Redis client initialized")

	if err := redis.InitStorage(); err != nil {
		log.Fatal(err)
	}
	defer redis.CloseStorageClient()
	log.Println("Redis storage initialized")

	seedURLs := []string{"https://www.naiin.com/"}

	go func() {
		err := crawler.Crawl(seedURLs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}

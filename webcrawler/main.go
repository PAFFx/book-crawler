package main

import (
	"context"
	"log"

	"book-search/webcrawler/services/redis"
)

func main() {
	ctx := context.Background()

	// Init services
	err := redis.InitRedisClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis client initialized")
}

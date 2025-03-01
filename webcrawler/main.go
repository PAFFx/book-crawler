package main

import (
	"log"

	"book-search/webcrawler/crawler"
	"book-search/webcrawler/services/redis"
)

func main() {
	// Init services
	if err := redis.InitStorage(); err != nil {
		log.Fatal(err)
	}
	log.Println("Redis storage backend initialized")

	seedURLs := []string{"https://www.chulabook.com"}

	err := crawler.Crawl(seedURLs)
	if err != nil {
		log.Fatal(err)
	}
}

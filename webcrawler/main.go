package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"book-search/webcrawler/crawler"
	"book-search/webcrawler/services/redis"
	"book-search/webcrawler/utils"
)

func main() {
	// Init services
	cleanupManager := utils.GetCleanupManager()
	cleanupManager.Add(func() {
		redis.CloseStorageClient()
	})

	// Run cleanup when recieve SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		cleanupManager.RunAll()
		os.Exit(0)
	}()

	seedURLs := []string{
		"https://www.chulabook.com",
		//"https://www.naiin.com",
	}

	err := crawler.Crawl(context.Background(), seedURLs)
	if err != nil {
		log.Fatal(err)
	}
}

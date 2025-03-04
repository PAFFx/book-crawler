package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"book-search/webcrawler/crawler"
	"book-search/webcrawler/services/database"
	"book-search/webcrawler/services/htmlStore"
	"book-search/webcrawler/services/storage"
	"book-search/webcrawler/utils"

	"github.com/gocolly/redisstorage"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func main() {
	cleanupManager := utils.GetCleanupManager()
	defer cleanupManager.RunAll()

	// Run cleanup when recieve SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		cleanupManager.RunAll()
		os.Exit(0)
	}()

	storageClient, htmlStoreClient, dbClient, err := initServices()
	if err != nil {
		log.Fatal(err)
	}

	seedURLs := []string{
		"https://www.chulabook.com",
		"https://www.naiin.com",
		"https://www.booktopia.com.au/books/fiction/cF-p1.html",
	}

	err = crawler.Crawl(context.Background(), storageClient, htmlStoreClient, dbClient, seedURLs)
	if err != nil {
		log.Fatal(err)
	}
}

func initServices() (*redisstorage.Storage, *minio.Client, *gorm.DB, error) {
	cleanupManager := utils.GetCleanupManager()

	// Init services
	storageClient, err := storage.GetStorage()
	if err != nil {
		log.Fatal(err)
	}
	cleanupManager.Add(func() { storage.CloseStorageClient(storageClient) })
	log.Println("Storage client init")

	htmlStoreClient, err := htmlStore.GetMinioClient() // no need to cleanup
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTML store client init")

	dbClient, err := database.GetDBClient()
	if err != nil {
		log.Fatal(err)
	}
	cleanupManager.Add(func() { database.CloseDBClient(dbClient) })
	log.Println("DB client init")

	return storageClient, htmlStoreClient, dbClient, nil

}

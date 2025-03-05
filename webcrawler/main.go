package main

import (
	"book-search/webcrawler/crawler"
	"book-search/webcrawler/services/database"
	"book-search/webcrawler/services/htmlStore"
	"book-search/webcrawler/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// main is the entry point for the web crawler application.
// It initializes services, sets up signal handling, and launches domain-specific crawlers.
func main() {
	cleanupManager := utils.GetCleanupManager()
	defer cleanupManager.RunAll()

	// Run cleanup when receive SIGINT or SIGTERM
	setupSignalHandling(cleanupManager)

	// Default seed URLs map - each domain has its own seed URLs
	seedURLMap := map[string][]string{
		"www.chulabook.com": {"https://www.chulabook.com"},
		"www.naiin.com":     {"https://www.naiin.com"},
		"www.booktopia.com.au": {
			"https://www.booktopia.com.au/books/fiction/cF-p1.html",
			"https://www.booktopia.com.au/ebooks/fiction/cF-p1-e.html",
			"https://www.booktopia.com.au/books/non-fiction/cN-p1.html",
			"https://www.booktopia.com.au/ebooks/non-fiction/cN-p1-e.html",
			"https://www.booktopia.com.au/books/text-books/higher-education-vocational-textbooks/cXA-p1.html",
			"https://www.booktopia.com.au/ebooks/non-fiction/accounting-finance/l101082-p1-e.html?cID=KF&sorter=bestsellers-dsc",
		},
	}

	// Initialize shared services
	htmlStoreClient, dbClient, err := initSharedServices()
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	log.Printf("Starting crawlers for each domain at %s...", startTime.Format(time.RFC3339))

	// Launch crawlers for each domain in the seed URLs map
	// This call blocks until all crawlers are finished
	crawler.LaunchCrawlers(seedURLMap, htmlStoreClient, dbClient)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("All crawlers have completed at %s (took %s)", endTime.Format(time.RFC3339), duration)
}

// setupSignalHandling configures signals to gracefully shut down the application.
// It listens for SIGINT and SIGTERM signals and runs cleanup handlers when received.
// Parameters:
//   - cleanupManager: Manager for running cleanup functions
func setupSignalHandling(cleanupManager *utils.CleanupManager) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		cleanupManager.RunAll()
		os.Exit(0)
	}()
}

// initSharedServices initializes services that are shared across all crawlers.
// These include HTML storage and database clients.
// Returns:
//   - *minio.Client: Client for HTML storage
//   - *gorm.DB: Database client
//   - error: Any error that occurred during initialization
func initSharedServices() (*minio.Client, *gorm.DB, error) {
	cleanupManager := utils.GetCleanupManager()

	htmlStoreClient, err := htmlStore.GetMinioClient() // no need to cleanup
	if err != nil {
		return nil, nil, err
	}
	log.Println("HTML store client init")

	dbClient, err := database.GetDBClient()
	if err != nil {
		return nil, nil, err
	}
	cleanupManager.Add(func() { database.CloseDBClient(dbClient) })
	log.Println("DB client init")

	return htmlStoreClient, dbClient, nil
}

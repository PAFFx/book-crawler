package crawler

import (
	"book-search/webcrawler/services/storage"
	"book-search/webcrawler/utils"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// LaunchCrawlers starts a separate crawler for each domain in the seed URL map
// Parameters:
//   - seedURLMap: Map of domains to their seed URLs
//   - htmlStoreClient: Client for storing HTML content
//   - dbClient: Database client for storing crawled data
func LaunchCrawlers(seedURLMap map[string][]string, htmlStoreClient *minio.Client, dbClient *gorm.DB) {
	var wg sync.WaitGroup
	activeCrawlers := 0

	// Get the progress manager and start periodic logging
	progressManager := GetProgressManager()
	stopLogging := progressManager.StartPeriodicLogging(5 * time.Second)
	defer stopLogging()

	utils.GetCleanupManager().Add(func() {
		stopLogging()
	})

	// Launch a dedicated crawler for each domain
	for domain, seedURLs := range seedURLMap {
		// Skip if no seeds
		if len(seedURLs) == 0 {
			log.Printf("No seed URLs for domain %s, skipping", domain)
			continue
		}

		activeCrawlers++
		wg.Add(1)

		// Create the Redis prefix for this domain
		redisPrefix := "webcrawler:" + domain

		// Launch crawler in a goroutine - this ensures they run in parallel
		go func(domain string, seedURLs []string, prefix string) {
			defer wg.Done()
			log.Printf("Starting crawler for domain: %s", domain)

			if err := launchSingleCrawler(domain, seedURLs, prefix, htmlStoreClient, dbClient); err != nil {
				log.Printf("Error with crawler for domain %s: %v", domain, err)
			} else {
				log.Printf("Crawler for domain %s completed successfully", domain)
			}
		}(domain, seedURLs, redisPrefix)
	}

	if activeCrawlers == 0 {
		log.Println("No active crawlers to run")
		return
	}

	log.Printf("Started %d crawlers in parallel", activeCrawlers)

	// Wait for all crawlers to finish
	wg.Wait()
	log.Println("All crawlers finished")
}

// launchSingleCrawler sets up and starts a crawler for a specific domain
// Parameters:
//   - domain: Domain to crawl
//   - seedURLs: Initial URLs to start crawling from
//   - redisPrefix: Prefix for Redis storage keys
//   - htmlStoreClient: Client for storing HTML content
//   - dbClient: Database client for storing crawled data
//
// Returns:
//   - error: Any error that occurred during setup or crawling
func launchSingleCrawler(domain string, seedURLs []string, redisPrefix string,
	htmlStoreClient *minio.Client, dbClient *gorm.DB) error {
	// Initialize storage with domain-specific prefix
	storageClient, err := storage.GetStorage(redisPrefix)
	if err != nil {
		return fmt.Errorf("failed to initialize storage for domain %s: %w", domain, err)
	}
	defer storage.CloseStorageClient(storageClient)

	// Run the crawler with only this domain in allowed domains
	err = Crawl(
		context.Background(),
		storageClient,
		htmlStoreClient,
		dbClient,
		seedURLs,
		[]string{domain},
	)
	if err != nil {
		return fmt.Errorf("crawler for domain %s failed: %w", domain, err)
	}

	return nil
}

package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	"book-search/webcrawler/config"
	"book-search/webcrawler/extracter"
	"book-search/webcrawler/services/database"
	"book-search/webcrawler/services/htmlStore"
	"book-search/webcrawler/utils"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func Crawl(ctx context.Context, storageClient *redisstorage.Storage, htmlStoreClient *minio.Client, dbClient *gorm.DB, seedURLs []string, allowedDomains []string) error {
	if len(allowedDomains) == 0 {
		return fmt.Errorf("no allowed domains specified")
	}

	domain := allowedDomains[0]
	// Get a progress tracker for this domain
	tracker := GetProgressManager().GetTracker(domain)

	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	// Create a collector with async mode enabled
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(allowedDomains...),
	)

	c.SetRequestTimeout(30 * time.Second)

	// Set domain-specific limits
	c.Limit(&colly.LimitRule{
		DomainGlob:  domain,
		Parallelism: 6, // Allow more parallelism since we're focused on one domain
		RandomDelay: 5 * time.Second,
		Delay:       1 * time.Second,
	})

	err = c.SetStorage(storageClient)
	if err != nil {
		return err
	}

	// Set up the queue
	q, err := queue.New(env.CrawlerThreads, storageClient)
	if err != nil {
		return err
	}

	// Set up request handlers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", config.GetRandomUserAgents())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Only follow links to the same domain (we're focusing on one domain per crawler)
		if e.Request.URL.Host == domain {
			err = q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
			if err != nil {
				log.Println("Error adding URL to queue:", err)
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		// Get a URL from the response
		urlStr := r.Request.URL.String()

		// Only track non-redirects (status 200)
		if r.StatusCode == 200 {
			// Track this page visit
			tracker.LogVisit(urlStr)

			e := extracter.GetExtracter(r.Request.URL.Host)
			if e != nil && e.IsValidBookPage(urlStr, string(r.Body)) {
				contentHash := utils.GenerateContentHash(string(r.Body))

				exists, err := database.CheckBookExists(dbClient, contentHash)
				if err != nil {
					log.Println("Error checking if book exists:", err)
				}

				if !exists {
					// Log book extraction in our tracker
					tracker.LogExtraction(urlStr)

					bookWithAuthors, err := e.Extract(string(r.Body))
					if err != nil {
						log.Println("Error extracting book:", err)
					}

					// Set HTMLHash for the book
					bookWithAuthors.Book.HTMLHash = contentHash

					err = database.StoreBookWithAuthors(dbClient, bookWithAuthors)
					if err != nil {
						log.Println("Error storing book with authors:", err)
					}

					err = htmlStore.StoreHTML(ctx, htmlStoreClient, string(r.Body), contentHash)
					if err != nil {
						log.Printf("Error storing HTML (URL: %s): %v\n", urlStr, err)
					}
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode != 404 && r.StatusCode != 500 {
			log.Println("Error visiting", r.Request.URL.String(), "with status", r.StatusCode)
		}
	})

	// Add the seed URLs to the queue
	for _, url := range seedURLs {
		if err = q.AddURL(url); err != nil {
			return err
		}
	}

	// Process the queue until it's empty
	for {
		// Run the queue with the async collector
		if err = q.Run(c); err != nil {
			return err
		}

		// Wait for all pending requests to complete before checking if queue is empty
		c.Wait()

		if q.IsEmpty() {
			break
		}

		// Small sleep to avoid CPU spinning
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// Helper function to truncate strings to a maximum length
func truncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength] + "..."
	}
	return str
}

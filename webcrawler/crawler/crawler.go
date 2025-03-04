package crawler

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"gorm.io/gorm"

	"book-search/webcrawler/config"
	"book-search/webcrawler/extracter"
	"book-search/webcrawler/services/database"
	"book-search/webcrawler/services/htmlStore"
	"book-search/webcrawler/utils"

	"github.com/minio/minio-go/v7"
)

func Crawl(ctx context.Context, storageClient *redisstorage.Storage, htmlStoreClient *minio.Client, dbClient *gorm.DB, seedURLs []string) error {
	cleanupManager := utils.GetCleanupManager()

	// get and check env
	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	bar := NewCrawlerProgressBar()
	cleanupManager.Add(func() {
		bar.Finish()
	})

	c := colly.NewCollector(
		colly.AllowedDomains(config.GetAllowedDomains()...),
		colly.Async(true),
	)

	c.SetRequestTimeout(30 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 6, // Concurrent request limit is not the same as the number of crawler consumer threads
		RandomDelay: 5 * time.Second,
		Delay:       1 * time.Second,
	})

	err = c.SetStorage(storageClient)
	if err != nil {
		return err
	}

	q, err := queue.New(env.CrawlerThreads, storageClient)
	if err != nil {
		return err
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", config.GetRandomUserAgents())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if slices.Contains(config.GetAllowedDomains(), e.Request.URL.Host) {
			err = q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
			if err != nil {
				log.Println("Error adding URL to queue:", err)
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			e := extracter.GetExtracter(r.Request.URL.Host)
			if e != nil && e.IsValidBookPage(r.Request.URL.String(), string(r.Body)) {

				contentHash := utils.GenerateContentHash(string(r.Body))

				exists, err := database.CheckBookExists(dbClient, contentHash)
				if err != nil {
					log.Println("Error checking if book exists:", err)
				}
				if !exists {

					const maxURLLength = 50
					bar.Describe("Crawled site: " + truncateString(r.Request.URL.String(), maxURLLength))
					_ = bar.Add(1)

					book, err := e.Extract(string(r.Body))
					if err != nil {
						log.Println("Error extracting book:", err)
					}

					err = database.StoreBook(dbClient, book)
					if err != nil {
						log.Println("Error storing book:", err)
					}

					err = htmlStore.StoreHTML(ctx, htmlStoreClient, string(r.Body), contentHash)
					if err != nil {
						log.Printf("Error storing HTML (URL: %s): %v\n", r.Request.URL.String(), err)
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

	for _, url := range seedURLs {
		if err = q.AddURL(url); err != nil {
			return err
		}
	}

	for {
		// Handle the queue with async requests,
		// wait for all requests to complete and check if the queue is empty
		if err = q.Run(c); err != nil {
			return err
		}
		c.Wait()
		if q.IsEmpty() {
			log.Println("Crawl completed, queue is empty")
			break
		}
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

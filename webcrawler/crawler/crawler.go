package crawler

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"

	"book-search/webcrawler/config"
	"book-search/webcrawler/extracter"
	"book-search/webcrawler/services/minio"
	"book-search/webcrawler/services/redis"
	"book-search/webcrawler/utils"
)

var allowedDomains = []string{"www.naiin.com", "www.chulabook.com", "www.amazon.com"}

func Crawl(ctx context.Context, seedURLs []string) error {
	// get and check env
	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	storage, err := redis.GetStorage()
	if err != nil {
		return err
	}

	bar := NewCrawlerProgressBar()

	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		colly.Async(true),
	)

	c.SetRequestTimeout(10 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 6, // Concurrent request limit is not the same as the number of crawler consumer threads
		RandomDelay: 5 * time.Second,
		Delay:       1 * time.Second,
	})

	err = c.SetStorage(storage)
	if err != nil {
		return err
	}

	q, err := queue.New(env.CrawlerThreads, storage)
	if err != nil {
		return err
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", config.GetRandomUserAgents())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if slices.Contains(allowedDomains, e.Request.URL.Host) {
			err = q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
			if err != nil {
				log.Println("Error adding URL to queue:", err)
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			_ = bar.Add(1)
			bar.Describe("Crawled site: " + r.Request.URL.String())

			// TODO: Check if the content is unique by comparing page hash first to save bandwidth
			e := extracter.GetExtracter(r.Request.URL.Host)
			if e != nil && e.IsValidBookPage(r.Request.URL.String(), string(r.Body)) {
				_, err := e.Extract(string(r.Body)) // TODO: Store the book data
				if err != nil {
					log.Println("Error extracting book:", err)
				}

				_, err = minio.StoreHTML(ctx, string(r.Body))
				if err != nil {
					log.Printf("Error storing HTML (URL: %s): %v\n", r.Request.URL.String(), err)
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error visiting", r.Request.URL.String(), "with status", r.StatusCode)
	})

	for _, url := range seedURLs {
		if err = q.AddURL(url); err != nil {
			return err
		}
	}

	cleanupManager := utils.GetCleanupManager()
	cleanupManager.Add(func() {
		// finish bar and
		bar.Finish()
	})

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

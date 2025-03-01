package crawler

import (
	"log"
	"slices"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"

	"book-search/webcrawler/config"
	"book-search/webcrawler/services/redis"
)

var allowedDomains = []string{"www.naiin.com", "www.chulabook.com", "www.amazon.com"}

func Crawl(seedURLs []string) error {
	// get and check env
	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	storage := redis.GetStorage()

	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		colly.UserAgent(env.CrawlerUserAgent),
		colly.Async(true),
	)

	err = c.SetStorage(storage)
	if err != nil {
		return err
	}

	q, err := queue.New(env.CrawlerThreads, storage)
	if err != nil {
		return err
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if slices.Contains(allowedDomains, e.Request.URL.Host) {
			log.Println("Adding URL to queue:", e.Request.AbsoluteURL(e.Attr("href")))
			err = q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
			if err != nil {
				log.Println("Error adding URL to queue:", err)
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL.String(), "with status", r.StatusCode)
	})

	for _, url := range seedURLs {
		if err = q.AddURL(url); err != nil {
			return err
		}
	}

	if err = q.Run(c); err != nil {
		return err
	}

	return nil
}

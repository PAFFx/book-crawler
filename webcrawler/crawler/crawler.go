package crawler

import (
	"log"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"

	"book-search/webcrawler/config"
)

var env *config.EnvVariables

func Crawl(url string, htmlStream chan<- string) {
	if env == nil {
		var err error
		env, err = config.GetEnv()
		if err != nil {
			log.Fatalf("Error getting environment variables: %v", err)
		}

		if env.CrawlerThreads <= 0 {
			log.Fatalf("Crawler threads must be greater than 0")
		}

		if env.CrawlerUserAgent == "" {
			log.Fatalf("Crawler user agent must be set")
		}
	}

	q, _ := queue.New(env.CrawlerThreads, &queue.InMemoryQueueStorage{MaxSize: 10000})
	c := colly.NewCollector(
		//colly.AllowedDomains("naiin", "chula", "amazon"),
		colly.UserAgent(env.CrawlerUserAgent),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//fmt.Println(e.Attr("href"))
		q.AddURL(e.Attr("href"))
	})

	c.OnResponse(func(r *colly.Response) {
		htmlStream <- string(r.Body) // Send the HTML content to the channel
	})

	q.AddURL(url)
	q.Run(c)
	close(htmlStream) // Close the channel when done
}

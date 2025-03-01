package main

import (
	"log"

	"book-search/webcrawler/crawler"
)

func main() {
	// Init services
	seedURLs := []string{
		"https://www.chulabook.com",
		"https://www.naiin.com",
	}

	err := crawler.Crawl(seedURLs)
	if err != nil {
		log.Fatal(err)
	}
}

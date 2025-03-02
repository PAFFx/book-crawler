package config

import (
	"math/rand"
)

var userAgents = []string{
	"BookBot/1.0",
	"BookScraper/1.0 (Book metadata collection bot",
	"BookCrawler/1.0 (Literature indexing service)",
	"LibraryCrawler/1.0",
	"BookIndexBot/2.0 (Book metadata harvester)",
	"MetadataBot/1.0 (Book information collection service)",
	"LiteratureBot/1.0",
}

func GetUserAgents() []string {
	return userAgents
}

func GetRandomUserAgents() string {
	return userAgents[rand.Intn(len(userAgents))]
}

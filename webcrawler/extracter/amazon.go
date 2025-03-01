package extracter

import (
	"book-search/webcrawler/models"
)

type AmazonExtracter struct{}

func (a AmazonExtracter) IsValidBookPage(url string, html string) bool {
	// Implement logic to check if the HTML is a valid Amazon book page
	return false
}

func (a AmazonExtracter) Extract(html string) (*models.Book, error) {
	// Implement logic to extract book information from Amazon HTML
	return &models.Book{}, nil
}

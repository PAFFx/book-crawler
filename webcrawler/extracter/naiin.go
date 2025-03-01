package extracter

import (
	"book-search/webcrawler/models"
)

type NaiinExtracter struct{}

func (n NaiinExtracter) IsValidBookPage(url string, html string) bool {
	// Implement logic to check if the HTML is a valid Naiin book page
	return false
}

func (n NaiinExtracter) Extract(html string) (*models.Book, error) {
	// Implement logic to extract book information from Naiin HTML
	return &models.Book{}, nil
}

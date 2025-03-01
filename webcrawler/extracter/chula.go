package extracter

import "book-search/webcrawler/models"

type ChulaExtracter struct{}


func (c ChulaExtracter) IsValidBookPage(url string, html string) bool {
	// Implement logic to check if the HTML is a valid Chula book page
	return false
}

func (n ChulaExtracter) Extract(html string) (*models.Book, error) {
	// Implement logic to extract book information from Naiin HTML
	return &models.Book{}, nil
}

package extracter

import (
	"book-search/webcrawler/models"
	"strings"
)

type Extracter interface {
	// IsValidBookPage checks if the html is a valid book page. Valid book pages are pages that contain book information.
	IsValidBookPage(url string, html string) bool

	// ExtractBookInfo extracts the book information from the html and returns a Book struct.
	Extract(html string) (*models.Book, error)
}

func GetExtracter(hostUrl string) Extracter {
	if strings.Contains(hostUrl, "amazon.com") {
		return &AmazonExtracter{}
	}

	if strings.Contains(hostUrl, "chulabook.com") {
		return &ChulaExtracter{}
	}
	return nil
}

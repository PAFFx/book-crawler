package extracter

import (
	"book-search/webcrawler/models"
	"strings"
)

type Extracter interface {
	// IsValidBookPage checks if the html is a valid book page. Valid book pages are pages that contain book information.
	IsValidBookPage(url string, html string) bool

	// Extract extracts the book and author information from the html and returns a BookWithAuthors struct.
	Extract(html string) (*models.BookWithAuthors, error)
}

func GetExtracter(hostUrl string) Extracter {
	if strings.Contains(hostUrl, "naiin.com") {
		return &NaiinExtracter{}
	}

	if strings.Contains(hostUrl, "chulabook.com") {
		return &ChulaExtracter{}
	}

	if strings.Contains(hostUrl, "booktopia.com.au") {
		return &BooktopiaExtracter{}
	}
	return nil
}

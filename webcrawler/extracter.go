package main

import "book-search/webcrawler/models"

type Extracter interface {
	// IsValidBookPage checks if the html is a valid book page. Valid book pages are pages that contain book information.
	IsValidBookPage(html string) bool

	// ExtractBookInfo extracts the book information from the html and returns a Book struct.
	Extract(html string) (*models.Book, error)
}

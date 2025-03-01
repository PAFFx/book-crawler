package models

import "net/url"

type Book struct {
	ProductURL  *url.URL
	ImageURL    *url.URL
	Title       string
	Authors     []string
	ISBN        string
	Description string
}

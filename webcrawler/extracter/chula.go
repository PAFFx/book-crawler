package extracter

import (
	"book-search/webcrawler/models"
	"book-search/webcrawler/utils"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ChulaExtracter struct {
}

func (c ChulaExtracter) IsValidBookPage(url string, html string) bool {
	// Implement logic to check if the HTML is a valid Chula book page
	if url != "" && strings.HasPrefix(url, "https://www.chulabook.com/") {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return false
		}
		description := strings.TrimSpace(doc.Find("h2:contains('รายละเอียดสินค้า')").Next().Text())
		authors := strings.TrimSpace(doc.Find(".detail-author").Text())
		authors = strings.Replace(authors, "ผู้แต่ง :", "", -1)

		if description != "" && authors != "" {
			return true
		}
		return false
	}
	return false
}

func (c ChulaExtracter) Extract(html string) (*models.BookWithAuthors, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Extract title
	title := doc.Find("title").Text()
	description := strings.TrimSpace(doc.Find("h2:contains('รายละเอียดสินค้า')").Next().Text())

	// Extract authors
	authorsText := strings.TrimSpace(doc.Find(".detail-author").Text())
	authorsText = strings.Replace(authorsText, "ผู้แต่ง :", "", -1)
	authorsText = strings.TrimSpace(authorsText)

	var authors []string
	if authorsText != "" {
		// Split by comma or other separators if needed
		authors = []string{authorsText} // Default to single author if no splitting needed
	}

	// Extract ISBN
	isbn := ""
	isbnText := doc.Find("p:contains('ISBN :')").Text()
	if isbnText != "" {
		isbn = isbnText
	}

	// Extract product URL
	var productURL *url.URL
	P_URL, exists := doc.Find(`meta[property="og:url"]`).Attr("content")
	if exists {
		parsedProductURL, err := url.Parse(P_URL)
		if err != nil {
			return nil, err
		}
		productURL = parsedProductURL
	}

	// Extract image URL
	var imageURL *url.URL
	Img_URL, exists := doc.Find(`meta[name="twitter:image"]`).Attr("content")
	if exists {
		parsedImageURL, err := url.Parse(Img_URL)
		if err != nil {
			return nil, err
		}
		imageURL = parsedImageURL
	}

	contentHash := utils.GenerateContentHash(html)

	book := &models.Book{
		HTMLHash:    contentHash,
		Title:       title, //NOTE GORM usually guess the relationship as has one if override foreign key name already exists in owner’s type, we need to specify references in the belongs to relationship.
		ISBN:        isbn,
		Description: description,
	}

	if productURL != nil {
		book.URL = productURL.String()
	}

	if imageURL != nil {
		book.ImageURL = imageURL.String()
	}

	return &models.BookWithAuthors{
		Book:    book,
		Authors: authors,
	}, nil
}

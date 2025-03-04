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

func (c ChulaExtracter) Extract(html string) (*models.Book, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Extract title
	title := doc.Find("title").Text()
	description := strings.TrimSpace(doc.Find("h2:contains('รายละเอียดสินค้า')").Next().Text())

	// Extract authors
	authors := strings.TrimSpace(doc.Find(".detail-author").Text())
	authors = strings.Replace(authors, "ผู้แต่ง :", "", -1)

	// Extract ISBN
	isbn := doc.Find("p:contains('ISBN :')").Text()

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

	// Create a new Book instance
	book := &models.Book{
		HTMLHash: contentHash,
		URL:      productURL.String(),
		ImageURL: imageURL.String(),
		Title:    title,
		//		Authors:     []string{authors},
		ISBN:        isbn,
		Description: description,
	}

	return book, nil
}

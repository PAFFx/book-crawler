package extracter

import (
	"github.com/PuerkitoBio/goquery"
	"book-search/webcrawler/models"
	"strings"
	"net/url"
	"fmt"
)

type ChulaExtracter struct {
	ProductURL *url.URL
	ImageURL   *url.URL
}

func (c ChulaExtracter) IsValidBookPage(url string, html string) bool {
	// Implement logic to check if the HTML is a valid Chula book page
	if url != "" && strings.HasPrefix(url, "https://www.chulabook.com/") {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return false
		}
		title := doc.Find(".product-name").Text()
		if title == "" {
			return false
		}
		fmt.Println("True")
		return true
	}
	fmt.Println("False")
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
	isbn := doc.Find("li:contains('ISBN')").Text()
	if strings.HasPrefix(isbn, "ISBN :") {
		isbn = strings.TrimSpace(strings.Replace(isbn, "ISBN :", "", -1))
	}

	// Extract product URL
	productURL, exists := doc.Find(`meta[property="og:url"]`).Attr("content")
	if exists {
		c.ProductURL, _ = url.Parse(productURL)
	}

	// Extract image URL
	imageURL, exists := doc.Find(`meta[name="twitter:image"]`).Attr("content")
	if exists {
		c.ImageURL, _ = url.Parse(imageURL)
	}

	// Create a new Book instance
	fmt.Println(c.ProductURL)
	fmt.Println(c.ImageURL)
	fmt.Println(title)
	fmt.Println(authors)
	fmt.Println(isbn)
	fmt.Println(description)
	book := &models.Book{
		Title:       title,
		Authors:     []string{authors},
		ISBN:        isbn,
		Description: description,
		ProductURL:  c.ProductURL,
		ImageURL:    c.ImageURL,
	}

	return book, nil
}

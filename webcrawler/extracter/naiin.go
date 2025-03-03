package extracter

import (
	"book-search/webcrawler/models"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var NAIIN_PRODUCT_TYPE_SELECTOR = "head > meta[property='og:type']" // content
var NAIIN_PRODUCT_URL_SELECTOR = "head > meta[property='og:url']"   // content
var NAIIN_IMAGE_URL_SELECTOR = "head > meta[property='og:image']"   // content
var NAIIN_TITLE_SELECTOR = ".bookdetail-container .title-topic"     // textContent
var NAIIN_BOOK_DETAIL_SELECTOR = ".bookdetail-container p"          // one that contains "ผู้เขียน:" -> split(",") -each> (replace("ผู้เขียน:", "") -> trim())
var NAIIN_ISBN_SELECTOR = "head > meta[property='book:isbn']"       // content
var NAIIN_DESCRIPTION_SELECTOR = ".book-decription"                 // textContent

type NaiinExtracter struct{}

func (n NaiinExtracter) IsValidBookPage(url string, html string) bool {
	matched, _ := regexp.MatchString(`https://www\.naiin\.com/product/detail/\d+`, url)
	if !matched {
		return false
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return false
	}

	productType := strings.TrimSpace(doc.Find(NAIIN_PRODUCT_TYPE_SELECTOR).First().AttrOr("content", ""))
	return strings.ToLower(productType) == "book"
}

func (n NaiinExtracter) Extract(html string) (*models.Book, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	productUrlStr := strings.TrimSpace(doc.Find(NAIIN_PRODUCT_URL_SELECTOR).First().AttrOr("content", ""))
	productUrl, err := url.Parse(productUrlStr)
	if err != nil {
		return nil, err
	}

	imageUrlStr := strings.TrimSpace(doc.Find(NAIIN_IMAGE_URL_SELECTOR).First().AttrOr("content", ""))
	imageUrl, err := url.Parse(imageUrlStr)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(doc.Find(NAIIN_TITLE_SELECTOR).First().Text())
	authors := []string{}
	doc.Find(NAIIN_BOOK_DETAIL_SELECTOR).Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "ผู้เขียน:") {
			for _, author := range strings.Split(s.Text(), ",") {
				author = strings.TrimSpace(strings.Replace(author, "ผู้เขียน:", "", 1))
				authors = append(authors, author)
			}
		}
	})

	isbn := strings.TrimSpace(doc.Find(NAIIN_ISBN_SELECTOR).First().AttrOr("content", ""))

	description := strings.TrimSpace(doc.Find(NAIIN_DESCRIPTION_SELECTOR).First().Text())

	return &models.Book{
		ProductURL:  productUrl,
		ImageURL:    imageUrl,
		Title:       title,
		Authors:     authors,
		ISBN:        isbn,
		Description: description,
	}, nil
}

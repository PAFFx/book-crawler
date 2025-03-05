package extracter

import (
	"book-search/webcrawler/models"
	"book-search/webcrawler/utils"
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

func (n NaiinExtracter) Extract(html string) (*models.BookWithAuthors, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Safe function to parse URLs with error handling
	safeParseURL := func(urlStr string) string {
		if urlStr == "" {
			return ""
		}

		// Remove any invalid control characters from URL
		cleanUrlStr := strings.Map(func(r rune) rune {
			if r < 32 || r == 127 { // ASCII control characters
				return -1 // Drop the character
			}
			return r
		}, urlStr)

		parsed, err := url.Parse(cleanUrlStr)
		if err != nil {
			// Log the error but return a safe empty string
			return ""
		}
		return parsed.String()
	}

	productUrlStr := strings.TrimSpace(doc.Find(NAIIN_PRODUCT_URL_SELECTOR).First().AttrOr("content", ""))
	imageUrlStr := strings.TrimSpace(doc.Find(NAIIN_IMAGE_URL_SELECTOR).First().AttrOr("content", ""))

	title := strings.TrimSpace(doc.Find(NAIIN_TITLE_SELECTOR).First().Text())
	authors := []string{}
	doc.Find(NAIIN_BOOK_DETAIL_SELECTOR).Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "ผู้เขียน:") {
			for _, author := range strings.Split(s.Text(), ",") {
				author = strings.TrimSpace(strings.Replace(author, "ผู้เขียน:", "", 1))
				if author != "" {
					authors = append(authors, author)
				}
			}
		}
	})

	isbn := strings.TrimSpace(doc.Find(NAIIN_ISBN_SELECTOR).First().AttrOr("content", ""))
	description := strings.TrimSpace(doc.Find(NAIIN_DESCRIPTION_SELECTOR).First().Text())

	// Create book with sanitized URLs
	book := &models.Book{
		HTMLHash:    utils.GenerateContentHash(html),
		URL:         safeParseURL(productUrlStr),
		ImageURL:    safeParseURL(imageUrlStr),
		Title:       title,
		ISBN:        isbn,
		Description: description,
	}

	return &models.BookWithAuthors{
		Book:    book,
		Authors: authors,
	}, nil
}

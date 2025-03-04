package extracter

import (
	"book-search/webcrawler/models"
	"book-search/webcrawler/utils"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var BOOKTOPIA_PRODUCT_TYPE_SELECTOR = "head > meta[property='og:type']"         // content
var BOOKTOPIA_PRODUCT_URL_SELECTOR = "head > meta[property='og:url']"           // content
var BOOKTOPIA_IMAGE_URL_SELECTOR = "head > meta[property='og:image']"           // content
var BOOKTOPIA_TITLE_SELECTOR = "#ProductDetails_d-product-info__rehyy h1"       // textContent
var BOOKTOPIA_PRODUCT_INFO_SELECTOR = "#ProductDetails_d-product-info__rehyy p" // one that contains "By:" -> replace("By:", "") -> split(",") -each> trim() -each> split(/\s+/).join("")
var BOOKTOPIA_DETAILS_SELECTOR = "#pdp-tabpanel-details p"                      // one that contains "ISBN:" -> replace("ISBN:", "") -> trim
var BOOKTOPIA_DESCRIPTION_SELECTOR = "#pdp-tabpanel-description"                // textContent

var reWhiteSpace, _ = regexp.Compile(`\s+`)

type BooktopiaExtracter struct{}

func (b BooktopiaExtracter) IsValidBookPage(url string, html string) bool {
	matched, _ := regexp.MatchString(`https://www\.booktopia\.com\.au/[^/]+/book/\d+\.html`, url)
	if !matched {
		return false
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return false
	}

	productType := strings.TrimSpace(doc.Find(BOOKTOPIA_PRODUCT_TYPE_SELECTOR).First().AttrOr("content", ""))
	return strings.ToLower(productType) == "book"
}

func (b BooktopiaExtracter) Extract(html string) (*models.Book, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	productUrlStr := strings.TrimSpace(doc.Find(BOOKTOPIA_PRODUCT_URL_SELECTOR).First().AttrOr("content", ""))
	productUrl, err := url.Parse(productUrlStr)
	if err != nil {
		return nil, err
	}

	imageUrlStr := strings.TrimSpace(doc.Find(BOOKTOPIA_IMAGE_URL_SELECTOR).First().AttrOr("content", ""))
	imageUrl, err := url.Parse(imageUrlStr)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(doc.Find(BOOKTOPIA_TITLE_SELECTOR).First().Text())

	authors := []string{}
	doc.Find(BOOKTOPIA_PRODUCT_INFO_SELECTOR).Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "By:") {
			ss := strings.Replace(s.Text(), "By:", "", 1)
			for _, author := range strings.Split(ss, ",") {
				author = strings.TrimSpace(author)
				authorNameParts := reWhiteSpace.Split(author, -1)
				author = strings.Join(authorNameParts, " ")
				authors = append(authors, author)
			}
		}
	})

	isbn := ""
	doc.Find(BOOKTOPIA_DETAILS_SELECTOR).Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "ISBN:") {
			isbn = strings.TrimSpace(strings.Replace(s.Text(), "ISBN:", "", 1))
		}
	})

	description := strings.TrimSpace(doc.Find(BOOKTOPIA_DESCRIPTION_SELECTOR).Text())

	contentHash := utils.GenerateContentHash(html)

	return &models.Book{
		HTMLHash: contentHash,
		URL:      productUrl.String(),
		ImageURL: imageUrl.String(),
		Title:    title,
		//Authors:     authors,
		ISBN:        isbn,
		Description: description,
	}, nil
}

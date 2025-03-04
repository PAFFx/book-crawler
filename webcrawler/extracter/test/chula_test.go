package extracter_test

import (
	"book-search/webcrawler/extracter"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestExtract(t *testing.T) {
	url := "https://www.chulabook.com/test-prep/208331"
	response, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to fetch URL: %v", err)
	}
	defer response.Body.Close()

	html, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	extracter := extracter.ChulaExtracter{}
	if !extracter.IsValidBookPage(url, string(html)) {
		fmt.Println("Not a valid book page")
	}
	bookWithAuthors, err := extracter.Extract(string(html))
	book := bookWithAuthors.Book

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if book.Title == "" {
		t.Errorf("Expected a title, got '%s'", book.Title)
	}
}

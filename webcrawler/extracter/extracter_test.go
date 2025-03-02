package extracter

import (
	"fmt"
	"io/ioutil"
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

	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	extracter := ChulaExtracter{}
	if !extracter.IsValidBookPage(url, string(html)) {
		fmt.Println("Not a valid book page")
	}
	book, err := extracter.Extract(string(html))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if book.Title == "" {
		t.Errorf("Expected a title, got '%s'", book.Title)
	}
}

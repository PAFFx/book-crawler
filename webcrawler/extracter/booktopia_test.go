package extracter

import (
	"book-search/webcrawler/config"
	// "encoding/json"
	// "fmt"
	"io"

	"net/http"
	"net/http/cookiejar"

	// "os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBooktopiaxtracter_IsValidBookPage(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid 1",
			url:  "https://www.booktopia.com.au/jurassic-park-michael-crichton/book/9780345538987.html",
			want: true,
		},
		{
			name: "Valid 2",
			url:  "https://www.booktopia.com.au/more-than-just-a-dog-simon-wooler/book/9780008707484.html",
			want: true,
		},
		{
			name: "Valid 3",
			url:  "https://www.booktopia.com.au/the-c-programming-language-brian-kernighan/book/9780131103627.html",
			want: true,
		},
		{
			name: "Invalid (main page)",
			url:  "https://www.booktopia.com.au/",
			want: false,
		},
		{
			name: "Invalid (fiction page)",
			url:  "https://www.booktopia.com.au/books/fiction/cF-p1.html",
			want: false,
		},
		{
			name: "Invalid (non-fiction page)",
			url:  "https://www.booktopia.com.au/books/non-fiction/cN-p1.html",
			want: false,
		},
		{
			name: "Invalid (textbook page)",
			url:  "https://www.booktopia.com.au/books/text-books/higher-education-vocational-textbooks/cXA-p1.html",
			want: false,
		},
		{
			name: "Invalid (textbook subject page)",
			url:  "https://www.booktopia.com.au/books/text-books/higher-education-vocational-textbooks/language-textbooks/cXAK-p1.html",
			want: false,
		},
		{
			name: "Invalid (stationery page)",
			url:  "https://www.booktopia.com.au/leuchtturm1917-notebook-medium-a5-hardcover-lined-black-leuchtturm1917/stationery/4004117258107.html",
			want: false,
		},
	}

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}
	b := BooktopiaExtracter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			req.Header.Set("User-Agent", config.GetRandomUserAgents())

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Failed to get URL: %v", err)
			}

			defer resp.Body.Close()
			html, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read body: %v", err)
			}

			// os.WriteFile(fmt.Sprintf("%s.html", tt.name), []byte(html), 0644)

			if got := b.IsValidBookPage(tt.url, string(html)); got != tt.want {
				t.Errorf("NaiinExtracter.IsValidBookPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooktopiaExtracter_Extract(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Book 1",
			url:  "https://www.booktopia.com.au/jurassic-park-michael-crichton/book/9780345538987.html",
		},
		{
			name: "Book 2",
			url:  "https://www.booktopia.com.au/more-than-just-a-dog-simon-wooler/book/9780008707484.html",
		},

		{
			name: "Book 3",
			url:  "https://www.booktopia.com.au/the-c-programming-language-brian-kernighan/book/9780131103627.html",
		},
	}

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}
	b := BooktopiaExtracter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			req.Header.Set("User-Agent", config.GetRandomUserAgents())

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Failed to get URL: %v", err)
			}

			defer resp.Body.Close()
			html, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read body: %v", err)
			}

			// os.WriteFile(fmt.Sprintf("%s.html", tt.name), []byte(html), 0644)

			book, errr := b.Extract(string(html))
			if errr != nil {
				t.Errorf("NaiinExtracter.Extract() error = %v", errr)
			}

			assert.NotEmpty(t, book.URL)
			assert.NotEmpty(t, book.ImageURL)
			assert.NotEmpty(t, book.Title)
			//assert.NotEmpty(t, book.Authors)
			assert.NotEmpty(t, book.ISBN)
			assert.NotEmpty(t, book.Description)

			// jsonData, _ := json.MarshalIndent(map[string]interface{}{
			// 	"product_url": book.ProductURL.String(),
			// 	"image_url":   book.ImageURL.String(),
			// 	"title":       book.Title,
			// 	"authors":     book.Authors,
			// 	"isbn":        book.ISBN,
			// 	"description": book.Description,
			// }, "", "  ")
			// os.WriteFile(fmt.Sprintf("%s.json", tt.name), jsonData, 0644)
		})
	}
}

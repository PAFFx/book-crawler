package extracter

import (
	// "encoding/json"
	// "fmt"
	"book-search/webcrawler/config"
	"io"

	"net/http"
	"net/http/cookiejar"

	// "os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNaiinExtracter_IsValidBookPage(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid 1",
			url:  "https://www.naiin.com/product/detail/639369",
			want: true,
		},
		{
			name: "Valid 2",
			url:  "https://www.naiin.com/product/detail/508064",
			want: true,
		},
		{
			name: "Invalid (main page)",
			url:  "https://www.naiin.com/books/",
			want: false,
		},
		{
			name: "Invalid (category page)",
			url:  "https://www.naiin.com/category?category_1_code=28&product_type_id=1",
			want: false,
		},
		{
			name: "Invalid (toy page)",
			url:  "https://www.naiin.com/product/detail/603593",
			want: false,
		},
	}

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}
	a := NaiinExtracter{}
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

			if got := a.IsValidBookPage(tt.url, string(html)); got != tt.want {
				t.Errorf("NaiinExtracter.IsValidBookPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNaiinExtracter_Extract(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Book 1",
			url:  "https://www.naiin.com/product/detail/639369",
		},
		{
			name: "Book 2",
			url:  "https://www.naiin.com/product/detail/508064",
		},

		{
			name: "Book 3",
			url:  "https://www.naiin.com/product/detail/485046",
		},
	}

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}
	a := NaiinExtracter{}
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

			book, errr := a.Extract(string(html))
			if errr != nil {
				t.Errorf("NaiinExtracter.Extract() error = %v", errr)
			}

			assert.NotEmpty(t, book.ProductURL.String())
			assert.NotEmpty(t, book.ImageURL.String())
			assert.NotEmpty(t, book.Title)
			assert.NotEmpty(t, book.Authors)
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

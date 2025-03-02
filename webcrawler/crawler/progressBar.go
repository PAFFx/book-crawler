package crawler

import (
	"github.com/schollz/progressbar/v3"
)

func NewCrawlerProgressBar() *progressbar.ProgressBar {
	// Create progress bar
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Crawling sites"),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}

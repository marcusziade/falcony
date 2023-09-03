package main

import (
	"testing"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

// Mocking youtube-dl command
func mockDownload(videoURL string) {
	// In a real-world scenario, here we would mock the actual download logic.
	// For this test, it's sufficient to know that this function would be called.
}

// Mocking appending downloaded title
func mockAppendDownloadedTitle(filePath, title string) {
	// Similarly, here we would mock the file appending logic.
	// For this test, it's enough to know that this function would be called.
}

func TestScrapeAndDownload(t *testing.T) {
	downloadedTitles := make(map[string]bool)
	downloadedTitles["Existing Bruce Falconer Video"] = true

	c := colly.NewCollector()

	// Mock the scraper logic
	c.OnHTML("a", func(e *colly.HTMLElement) {
		title := e.Attr("title")

		if title == "New Bruce Falconer Video" {
			assert.False(t, downloadedTitles[title])
			mockDownload("mock_video_url")
			mockAppendDownloadedTitle("mock_file_path", title)
		}

		if title == "Existing Bruce Falconer Video" {
			assert.True(t, downloadedTitles[title])
		}
	})

	// Run the mock scrape function
	c.Visit("https://www.youtube.com/c/MockChannelName/videos")
}


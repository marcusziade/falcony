package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

func readDownloadedTitles(filePath string) (map[string]bool, error) {
	titles := make(map[string]bool)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open downloaded titles file: %v", err)
		return titles, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		title := scanner.Text()
		titles[title] = true
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read from file: %v", err)
		return nil, err
	}

	return titles, nil
}

func appendDownloadedTitle(filePath, title string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open or create file: %v", err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(title + "\n"); err != nil {
		log.Printf("Failed to write to file: %v", err)
		return err
	}

	if err := writer.Flush(); err != nil {
		log.Printf("Failed to flush writer: %v", err)
		return err
	}

	return nil
}

func ScrapeAndDownload() {
	filePath := "app/videos/downloaded_titles.txt"
	if err := ensureFileExists(filePath); err != nil {
		log.Fatalf("Failed to ensure file exists: %v", err)
		return
	}

	downloadedTitles, err := readDownloadedTitles(filePath)
	if err != nil {
		log.Fatalf("Failed to read downloaded titles: %v", err)
	}

	c := colly.NewCollector()

	c.OnHTML("a", func(e *colly.HTMLElement) {
		title := e.Attr("title")

		if strings.Contains(title, "Bruce Falconer") && !downloadedTitles[title] {
			videoURL := "https://www.youtube.com" + e.Attr("href")
			cmd := exec.Command("youtube-dl", "-o", "/app/videos/%(title)s.%(ext)s", videoURL)
			if err := cmd.Run(); err != nil {
				log.Printf("Failed to download video: %v", err)
				return
			}

			appendDownloadedTitle("/app/videos/downloaded_titles.txt", title)

			log.Printf("Successfully downloaded and saved the video: %s", title)
		}
	})

	if err := c.Visit("https://www.youtube.com/@goodgrief3308/videos"); err != nil {
		log.Fatalf("Failed to visit the URL: %v", err)
	}
}

func ensureFileExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Failed to create directory: %v", err)
		return err
	}

	_, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}

	return nil
}

func main() {
	ScrapeAndDownload()
}

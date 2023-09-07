package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/kkdai/youtube/v2"
)

const channelURL = "https://www.youtube.com/user/@goodgrief3308/videos"
const historyFilePath = "videoHistory.json"

func getLatestVideoURL(ctx context.Context) (string, string, error) {
	log.Println("Navigating to channel...")
	var videoURL, thumbnailURL string
	err := chromedp.Run(ctx,
		chromedp.Navigate(channelURL),
		chromedp.WaitVisible(`a#thumbnail`, chromedp.ByQuery),
		chromedp.AttributeValue(`a#thumbnail img`, "src", &thumbnailURL, nil),
		chromedp.AttributeValue(`a#video-title-link`, "href", &videoURL, nil),
	)
	if videoURL == "" {
		return "", "", errors.New("could not find video URL")
	}
	return "https://www.youtube.com" + videoURL, thumbnailURL, err
}

func handleFatalError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func bestFormat(formats youtube.FormatList) *youtube.Format {
	log.Println("Determining best video format...")
	formats.Sort()
	return &formats[0]
}

func downloadVideo(ctx context.Context, videoURL string) {
	log.Println("Initiating video download...")

	client := youtube.Client{}

	videoID, err := youtube.ExtractVideoID(videoURL)
	handleFatalError(err, "extract video ID error")

	video, err := client.GetVideo(videoID)
	handleFatalError(err, "get video error")

	stream, _, err := client.GetStream(video, bestFormat(video.Formats))
	handleFatalError(err, "get video stream error")

	sanitizedTitle := strings.Map(func(r rune) rune {
		if r == os.PathSeparator || r == ':' {
			return '_'
		}
		return r
	}, video.Title)

	file, err := os.Create(sanitizedTitle + ".mp4")
	handleFatalError(err, "create file error")
	defer file.Close()

	log.Printf("Downloading video: %s...\n", video.Title)

	log.Println("Download started at:", time.Now().Format(time.RFC1123))

	_, err = io.Copy(file, stream)
	handleFatalError(err, "download video error")

	log.Println("Download completed at:", time.Now().Format(time.RFC1123))

	log.Printf("Downloaded: %s\n", video.Title)
}

func main() {
	log.Println("Starting application...")

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	videoURL, _, err := getLatestVideoURL(ctx)
	handleFatalError(err, "get latest video URL error")

	var history []string
	file, err := os.ReadFile(historyFilePath)
	if err == nil {
		log.Println("Reading history file...")
		json.Unmarshal(file, &history)
	} else if !os.IsNotExist(err) {
		log.Fatalf("Error reading history file: %v", err)
	}

	for _, v := range history {
		if v == videoURL {
			log.Println("No new videos to download.")
			return
		}
	}

	log.Println("Updating history file...")
	history = append(history, videoURL)
	file, err = json.Marshal(history)
	handleFatalError(err, "Error marshalling history data")

	err = os.WriteFile(historyFilePath, file, 0644)
	handleFatalError(err, "Error writing history file")

	downloadVideo(ctx, videoURL)
}

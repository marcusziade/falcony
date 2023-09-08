package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	log.Println("Starting program")
	ctx := context.Background()

	log.Println("Getting channel ID")
	channelID := "UCdlRNHC2zbCfTlpSiXKcDLg"

	log.Println("Getting API key")
	apiKey := getAPIKey()

	log.Println("Getting YouTube service")
	service := getYouTubeService(ctx, apiKey)

	log.Println("Getting latest video from channel")
	response := getLatestVideoFromChannel(service, channelID)

	if len(response.Items) == 0 {
		log.Println("No new videos found.")
		return
	}

	latestVideoDate := response.Items[0].Snippet.PublishedAt
	lastDownloadDate := readLastDownloadDate()

	if latestVideoDate == lastDownloadDate {
		log.Println("No new videos found.")
		return
	}

	videoID := response.Items[0].Id.VideoId

	log.Println("Downloading video using yt-dlp")
	err := downloadVideo(videoID)
	if err != nil {
		log.Fatalf("Error downloading video: %v", err)
	}

	log.Println("Video downloaded successfully: " + videoID + ".mp4")
	writeLastDownloadDate(latestVideoDate)
}

func downloadVideo(videoID string) error {
	cmd := exec.Command("yt-dlp", "-S", "ext:mp4:m4a", videoID)
	err := cmd.Run()
	return err
}

func getAPIKey() string {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		log.Fatalf("Error reading YOUTUBE_API_KEY environment variable")
	}
	return apiKey
}

func getYouTubeService(ctx context.Context, apiKey string) *youtube.Service {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	return service
}

func getLatestVideoFromChannel(service *youtube.Service, channelID string) *youtube.SearchListResponse {
	call := service.Search.List([]string{"id", "snippet"}).
		ChannelId(channelID).
		Order("date").
		MaxResults(1)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}
	return response
}

func readLastDownloadDate() string {
	data, err := os.ReadFile("last_download.txt")
	if err != nil {
		return ""
	}
	return string(data)
}

func writeLastDownloadDate(date string) {
	err := os.WriteFile("last_download.txt", []byte(date), 0644)
	if err != nil {
		log.Fatalf("Error writing to last_download.txt: %v", err)
	}
}

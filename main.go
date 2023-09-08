package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	ytdl "github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	log.Println("Starting program")
	ctx := context.Background()

	log.Println("Getting channel ID")
	channelID := getChannelID()

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
	client := ytdl.Client{}

	log.Println("Getting video details")
	video := getVideoDetails(ctx, client, videoID)

	log.Println("Getting video stream")
	stream := getVideoStream(ctx, client, video)

	log.Println("Reading stream into buffer")
	buf := readStreamIntoBuffer(stream)

	log.Println("Writing buffer to file")
	writeBufferToFile(video.Title, buf)

	log.Println("Video downloaded successfully: " + video.Title + ".mp4")
	writeLastDownloadDate(latestVideoDate)
}

func getChannelID() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the channel ID: ")
	channelID, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading channel ID: %v", err)
	}
	return channelID[:len(channelID)-1]
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

func getVideoDetails(ctx context.Context, client ytdl.Client, videoID string) *ytdl.Video {
	video, err := client.GetVideoContext(ctx, videoID)
	if err != nil {
		log.Fatalf("Error getting video details: %v", err)
	}
	return video
}

func getVideoStream(ctx context.Context, client ytdl.Client, video *ytdl.Video) io.ReadCloser {
	stream, _, err := client.GetStreamContext(ctx, video, &video.Formats[0])
	if err != nil {
		log.Fatalf("Error getting video stream: %v", err)
	}
	return stream
}

func readStreamIntoBuffer(stream io.ReadCloser) []byte {
	buf, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Fatalf("Error reading video stream: %v", err)
	}
	return buf
}

func writeBufferToFile(title string, buf []byte) {
	if _, err := os.Stat("videos/"); os.IsNotExist(err) {
		os.Mkdir("videos/", 0755)
	}

	file, err := os.Create("videos/" + title + ".mp4")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Error closing file: %v", err)
	}
}

func readLastDownloadDate() string {
	data, err := ioutil.ReadFile("last_download.txt")
	if err != nil {
		return ""
	}
	return string(data)
}

func writeLastDownloadDate(date string) {
	err := ioutil.WriteFile("last_download.txt", []byte(date), 0644)
	if err != nil {
		log.Fatalf("Error writing to last_download.txt: %v", err)
	}
}

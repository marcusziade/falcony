# YouTube Video Downloader 

## Description

This Go application fetches the latest video from a specified YouTube channel and downloads it using `yt-dlp`. It tracks the last downloaded video date to avoid duplicate downloads. The application can be dockerized and scheduled to run on a daily basis using a cron job.

## Prerequisites

Ensure you have the following installed on your local machine:

- Go (latest version)
- Docker
- ffmpeg
- curl

## Environment Variables

To run this application, you need to set the following environment variable:

- `YOUTUBE_API_KEY`: Your YouTube API key.

## Setup and Installation

1. Clone this repository to your local machine.
2. Navigate to the project directory.
3. Run `go mod download` to install the necessary Go packages.
4. Set up the YouTube API key as an environment variable.
5. Build the Go application using the command `go build -o main`.
6. Build the Docker image using the command `docker build -t youtube-video-downloader .`.
7. Run the application with `./main` or use Docker to run the application (see instructions below).

## How to Use

1. Set the channel ID in the `main()` function to your desired YouTube channel.
2. Run the application.
3. The application will fetch and download the latest video from the specified channel to the `videos` directory in the project.

## Running the Application Daily with Docker

1. First, ensure that the Docker daemon is running on your system.
2. Build the Docker image using the following command in the project directory:
   ```
   docker build -t youtube-video-downloader .
   ```
3. To run the container daily, you can set up a cron job. Open the crontab file with the command:
   ```
   crontab -e
   ```
4. Add the following line to run the container daily at a specific time (e.g., at 2 am):
   ```
   0 0,12 * * * docker run -e YOUTUBE_API_KEY='your-youtube-api-key' -d my-golang-app
   ```
5. Save and exit the crontab file. The docker container will now run daily at the specified time.
6. You can run it once now: `0 0,12 * * * docker run -e YOUTUBE_API_KEY='your-youtube-api-key' -d my-golang-app`. It will appear in Docker's desktop app where you can run it 
7. If you want to mount storage to the Docker automation:
```
docker run -e YOUTUBE_API_KEY='your-youtube-api-key' -v /yourpathyouwantToSave:/app/videos -d falcony
```

## Application Flow
### README.md

---

# YouTube Video Downloader 

## Description

This Go application fetches the latest video from a specified YouTube channel and downloads it using `yt-dlp`. It tracks the last downloaded video date to avoid duplicate downloads.

## Prerequisites

Ensure you have the following installed on your system:

- Go (latest version)
- ffmpeg
- curl

## Environment Variables

To run this application, you need to set the following environment variable:

- `YOUTUBE_API_KEY`: Your YouTube API key.

## Setup and Installation

1. Clone this repository to your local machine.
2. Navigate to the project directory.
3. Run `go mod download` to install the necessary Go packages.
4. Set up the YouTube API key as an environment variable.
5. Build the Go application using the command `go build -o main`.
6. Run the application with `./main`.

## How to Use

1. Set the channel ID in the `main()` function to your desired YouTube channel.
2. Run the application.
3. The application will fetch and download the latest video from the specified channel to the `videos` directory in the project.

## Application Flow

1. Checks and creates a `videos` directory if not exists.
2. Retrieves the channel ID and API key.
3. Initializes a YouTube service.
4. Fetches the latest video from the specified channel.
5. Compares the latest video's publish date with the last download date to prevent duplicate downloads.
6. Downloads the video using `yt-dlp` if a new video is found.
7. Updates the last download date.

## Error Handling

The application contains robust error handling to manage issues like:

- Failure to read the YouTube API key from environment variables.
- Errors during YouTube service initialization.
- Errors during API calls to fetch video details.
- Issues during video download.
- Failures to read or write the last download date.

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests.

## License

This project is open-source and available under the MIT License.

---

I hope this README.md serves your needs, Marcus. Let me know if there are any specific sections or details you'd like added or changed.

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests.

## License

This project is open-source and available under the MIT License.

---

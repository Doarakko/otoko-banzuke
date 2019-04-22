package youtube

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// NewYoutubeService hoge
func NewYoutubeService() *youtube.Service {
	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("YOUTUBE_API_KEY")},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	return service
}

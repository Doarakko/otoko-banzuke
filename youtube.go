package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type channel struct {
	id string
}

type video struct {
	id string
}

type comment struct {
	ID          string
	TextDisplay string
	AuthorID    string
	AuthorName  string
	AuthorURL   string
	ChannelID   string
	VideoID     string
	ParentID    string
	LikeCount   int16
	ReplyCount  int16
	CreateDate  time.Time
	UpdateDate  time.Time
}

func gormConnect() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err.Error())
	}
	return db
}

func youtubeService() *youtube.Service {
	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("YOTUBE_API_KEY")},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	return service
}

func getComments(VideoID string) []comment {
	service := youtubeService()
	call := service.CommentThreads.List("id,snippet").
		VideoId(VideoID).
		SearchTerms("草").
		Order("relevance").
		MaxResults(5)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	comments := []comment{}
	for _, item := range response.Items {
		id := item.Snippet.TopLevelComment.Id
		authorName := item.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		authorURL := item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl
		textDisplay := item.Snippet.TopLevelComment.Snippet.TextDisplay
		likeCount := int16(item.Snippet.TopLevelComment.Snippet.LikeCount)
		replyCount := int16(item.Snippet.TotalReplyCount)
		channelID := item.Snippet.ChannelId
		videoID := item.Snippet.VideoId

		comment := comment{
			ID:          id,
			VideoID:     videoID,
			ChannelID:   channelID,
			TextDisplay: textDisplay,
			AuthorName:  authorName,
			AuthorURL:   authorURL,
			LikeCount:   likeCount,
			ReplyCount:  replyCount,
		}
		comments = append(comments, comment)
	}

	return comments
}

func getVideoIDs(channelID string) {
}

func insertComment(comments []comment) {
	db := gormConnect()
	defer db.Close()

	for _, item := range comments {
		db.Create(&item)
	}
}

func main() {
	videoID := "Qs3sShlgKGk"
	comments := getComments(videoID)
	insertComment(comments)
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type channel struct {
	ChannelID    string
	Name         string
	Description  string
	ThumbnailURL string
	PlaylistID   string
	ViewCount    int64
	VideoCount   int32
	// CommentCount    int32
	SubscriberCount int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type video struct {
	VideoID      string
	Title        string
	Description  string
	ThumbnailURL string
	ViewCount    int64
	LikeCount    int32
	DislikeCount int32
	CommentCount int32
	ChannelID    string
	PublishedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type comment struct {
	CommentID   string
	TextDisplay string
	AuthorID    string
	AuthorName  string
	AuthorURL   string
	ChannelID   string
	VideoID     string
	ParentID    string
	LikeCount   int32
	ReplyCount  int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newGormConnect() *gorm.DB {
	//db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := gorm.Open("postgres", "postgres://uojcpxxbqfdtwe:32aaebe5b808ffc63f7753a98f16000479469eb57e7a85eed218feb2e0436463@ec2-23-23-173-30.compute-1.amazonaws.com:5432/d16518rirmmbsl")

	if err != nil {
		panic(err.Error())
	}
	return db
}

func newYoutubeService() *youtube.Service {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyD3BcAAS4dmpvUbVrT8K3U49d6PYxR8CTk"},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	return service
}

func getComments(VideoID string) []comment {
	service := newYoutubeService()
	call := service.CommentThreads.List("snippet").
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
		commentid := item.Snippet.TopLevelComment.Id
		authorName := item.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		authorURL := item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl
		textDisplay := item.Snippet.TopLevelComment.Snippet.TextDisplay
		likeCount := int32(item.Snippet.TopLevelComment.Snippet.LikeCount)
		replyCount := int32(item.Snippet.TotalReplyCount)
		channelID := item.Snippet.ChannelId
		videoID := item.Snippet.VideoId

		comment := comment{
			CommentID:   commentid,
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

func insertChannel(channel channel) {
	db := newGormConnect()
	defer db.Close()

	db.Create(&channel)
}

func insertVideos(videos []video) {
	db := newGormConnect()
	defer db.Close()

	for _, item := range videos {
		db.Create(&item)
	}
}

func insertComments(comments []comment) {
	db := newGormConnect()
	defer db.Close()

	for _, item := range comments {
		db.Create(&item)
	}
}

func selectChannel(channelID string) channel {
	db := newGormConnect()
	defer db.Close()

	channel := channel{}
	db.First(&channel, "channel_id=?", channelID)

	return channel
}

func getChannelInfo(channelID string) channel {
	service := newYoutubeService()
	call := service.Channels.List("snippet,contentDetails,statistics").
		Id(channelID).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	item := response.Items[0]

	name := item.Snippet.Title
	description := item.Snippet.Description
	thumbnailURL := item.Snippet.Thumbnails.High.Url
	playlistID := item.ContentDetails.RelatedPlaylists.Uploads
	viewCount := int64(item.Statistics.ViewCount)
	subscriberCount := int32(item.Statistics.SubscriberCount)
	videoCount := int32(item.Statistics.VideoCount)
	// commentCount := item.Statistics.CommentCount

	channel := channel{
		ChannelID:       channelID,
		Name:            name,
		Description:     description,
		ThumbnailURL:    thumbnailURL,
		PlaylistID:      playlistID,
		ViewCount:       viewCount,
		VideoCount:      videoCount,
		SubscriberCount: subscriberCount,
	}
	return channel
}

func getVideoInfo(videoID string) video {
	service := newYoutubeService()
	call := service.Videos.List("id,snippet,Statistics").
		Id(videoID).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	item := response.Items[0]

	title := item.Snippet.Title
	description := item.Snippet.Description
	thumbnailURL := item.Snippet.Thumbnails.High.Url
	viewCount := int64(item.Statistics.ViewCount)
	commentCount := int32(item.Statistics.CommentCount)
	likeCount := int32(item.Statistics.LikeCount)
	dislikeCount := int32(item.Statistics.DislikeCount)
	channelID := item.Snippet.ChannelId
	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		log.Fatalf("%v", err)
	}

	video := video{
		VideoID:      videoID,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		ViewCount:    viewCount,
		LikeCount:    likeCount,
		DislikeCount: dislikeCount,
		CommentCount: commentCount,
		ChannelID:    channelID,
		PublishedAt:  publishedAt,
	}

	return video
}

// TODO 日付指定
func getNewVideos(channelID string) []video {
	service := newYoutubeService()
	call := service.Search.List("id").
		Type("video").
		ChannelId(channelID).
		Order("date").
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []video{}
	for _, item := range response.Items {
		videoID := item.Id.VideoId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url
		channelID := item.Snippet.ChannelId
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("%v", err)
		}

		video := video{
			VideoID:      videoID,
			Title:        title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
			ChannelID:    channelID,
			PublishedAt:  publishedAt,
		}

		videos = append(videos, video)
	}
	return videos
}

func getAllVideos(playlistID string, pageToken string) []video {
	service := newYoutubeService()
	call := service.PlaylistItems.List("id,snippet,contentDetails").
		PlaylistId(playlistID).
		PageToken(pageToken).
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []video{}
	for _, item := range response.Items {
		videoID := item.ContentDetails.VideoId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url
		channelID := item.Snippet.ChannelId
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("%v", err)
		}

		video := video{
			VideoID:      videoID,
			Title:        title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
			ChannelID:    channelID,
			PublishedAt:  publishedAt,
		}

		videos = append(videos, video)
	}

	pageToken = response.NextPageToken
	if pageToken != "" {
		videos = append(videos, getAllVideos(playlistID, pageToken)...)
	}

	return videos
}

func main() {
	channelID := "UCENoC6MLc4pL-vehJyzSWmg"
	//insertChannel(getChannelInfo(channelID))
	channel := selectChannel(channelID)

	insertVideos(getAllVideos(channel.PlaylistID, ""))
}

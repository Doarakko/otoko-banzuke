package main

import (
	"log"
	"time"
)

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

func (v *video) insertVideo() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&v)
	log.Printf("Insert video: %v\n", v)
}

func (v *video) getVideoInfo() video {
	service := newYoutubeService()
	call := service.Videos.List("id,snippet,Statistics").
		Id(v.VideoID).
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

	*v = video{
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

	return *v
}

func (v *video) getComments() []comment {
	service := newYoutubeService()
	call := service.CommentThreads.List("snippet").
		VideoId(v.VideoID).
		Order("relevance").
		MaxResults(50)
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
	log.Printf("Get %v comments\n", len(comments))

	return comments
}

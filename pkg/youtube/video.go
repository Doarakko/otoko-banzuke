package youtube

import (
	"log"
	"time"
)

// Video gare
type Video struct {
	VideoID      string
	Title        string
	Description  string
	ThumbnailURL string
	ViewCount    int64
	CommentCount int32
	ChannelID    string
	PublishedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (v *Video) insertVideo() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&v)
	log.Printf("Insert video: %v\n", v)
}

func (v *Video) deleteVideo() {

}

// SetDetailInfo hoge
func (v *Video) SetDetailInfo() {
	service := newYoutubeService()
	call := service.Videos.List("id,snippet,Statistics").
		Id(v.VideoID).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	item := response.Items[0]

	v.Title = item.Snippet.Title
	v.Description = item.Snippet.Description
	v.ThumbnailURL = item.Snippet.Thumbnails.High.Url
	v.ViewCount = int64(item.Statistics.ViewCount)
	v.CommentCount = int32(item.Statistics.CommentCount)
	v.ChannelID = item.Snippet.ChannelId
	v.PublishedAt, err = time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (v *Video) getComments() []Comment {
	service := newYoutubeService()
	call := service.CommentThreads.List("snippet").
		VideoId(v.VideoID).
		Order("relevance").
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	comments := []Comment{}
	for _, item := range response.Items {
		commentid := item.Snippet.TopLevelComment.Id
		authorName := item.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		authorURL := item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl
		textDisplay := item.Snippet.TopLevelComment.Snippet.TextDisplay
		likeCount := int32(item.Snippet.TopLevelComment.Snippet.LikeCount)
		replyCount := int32(item.Snippet.TotalReplyCount)
		channelID := item.Snippet.ChannelId
		videoID := item.Snippet.VideoId

		comment := Comment{
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

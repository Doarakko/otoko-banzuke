package youtube

import (
	"log"
	"time"

	mydb "../database"
)

// Video gare
type Video struct {
	VideoID      string    `gorm:"column:video_id"`
	Title        string    `gorm:"column:title"`
	Description  string    `gorm:"column:description"`
	ThumbnailURL string    `gorm:"column:thumbnail_url"`
	ViewCount    int64     `gorm:"column:view_count"`
	CommentCount int32     `gorm:"column:comment_count"`
	PublishedAt  time.Time `gorm:"column:published_at"`
	CategoryID   string    `gorm:"column:category_id"`
	CategoryName string    `gorm:"column:category_name"`
	ChannelID    string    `gorm:"column:channel_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

// Insert video
func (v *Video) Insert() {
	db := mydb.NewGormConnect()
	defer db.Close()

	db.Create(&v)
	log.Printf("Insert video: %v\n", v)
}

// Delete video
func (v *Video) Delete() {

}

// SetDetailInfo hoge
func (v *Video) SetDetailInfo() {
	service := NewYoutubeService()
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
	v.CategoryID = item.Snippet.CategoryId
	v.setCategoryName()
	v.PublishedAt, err = time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (v *Video) setCategoryName() {
	service := NewYoutubeService()
	call := service.VideoCategories.List("id,snippet").
		Id(v.CategoryID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	item := response.Items[0]
	v.CategoryName = item.Snippet.Title
}

// GetComments hoge
func (v *Video) GetComments() []Comment {
	service := NewYoutubeService()
	call := service.CommentThreads.List("snippet").
		VideoId(v.VideoID).
		TextFormat("plainText").
		Order("relevance").
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	comments := []Comment{}
	for _, item := range response.Items {
		commentID := item.Snippet.TopLevelComment.Id
		authorName := item.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		authorURL := item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl
		textDisplay := item.Snippet.TopLevelComment.Snippet.TextDisplay
		likeCount := int32(item.Snippet.TopLevelComment.Snippet.LikeCount)
		replyCount := int32(item.Snippet.TotalReplyCount)

		comment := Comment{
			CommentID:   commentID,
			VideoID:     v.VideoID,
			ChannelID:   v.ChannelID,
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

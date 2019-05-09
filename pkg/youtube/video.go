package youtube

import (
	"log"
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	"google.golang.org/api/youtube/v3"
)

// Video struct
type Video struct {
	VideoID      string    `gorm:"column:video_id;primary_key"`
	Title        string    `gorm:"column:title;not null"`
	Description  string    `gorm:"column:description;not null"`
	ThumbnailURL string    `gorm:"column:thumbnail_url;not null"`
	ViewCount    int64     `gorm:"column:view_count;not null"`
	CommentCount int32     `gorm:"column:comment_count;not null"`
	PublishedAt  time.Time `gorm:"column:published_at;not null"`
	CategoryID   string    `gorm:"column:category_id;not null"`
	CategoryName string    `gorm:"column:category_name;not null"`
	ChannelID    string    `gorm:"column:channel_id;not null;index"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

// Exists if video exist return true
func (v *Video) Exists() bool {
	db := mydb.NewGormConnect()
	defer db.Close()

	result := db.First(&v, "video_id=?", v.VideoID)

	return result.Error == nil
}

// Insert video
func (v *Video) Insert() {
	db := mydb.NewGormConnect()
	defer db.Close()

	v.SetDetailInfo()
	db.Create(&v)
	log.Printf("Insert video: %v\n", v)
}

// Update video
func (v *Video) Update() {
	db := mydb.NewGormConnect()
	defer db.Close()

	v.SetDetailInfo()

	db.Model(&v).Updates(Video{
		Title:        v.Title,
		Description:  v.Description,
		ThumbnailURL: v.ThumbnailURL,
		ViewCount:    v.ViewCount,
		CommentCount: v.CommentCount,
		CategoryID:   v.CategoryID,
		CategoryName: v.CategoryName,
	})
	log.Printf("Update video: %v\n", v.VideoID)
}

// Delete video
func (v *Video) Delete() {

}

// SetDetailInfo ViewCount, CommentCount, CategoryID, CategoryName
func (v *Video) SetDetailInfo() {
	service := NewYoutubeService()
	call := service.Videos.List("snippet,Statistics").
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
	call := service.VideoCategories.List("snippet").
		Id(v.CategoryID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	item := response.Items[0]
	v.CategoryName = item.Snippet.Title
}

// GetComments get comment
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
	log.Printf("Get %v comments from %v\n", len(comments), v.VideoID)

	return comments
}

func newVideo(item youtube.SearchResult) Video {
	videoID := item.Id.VideoId
	title := item.Snippet.Title
	description := item.Snippet.Description
	thumbnailURL := item.Snippet.Thumbnails.High.Url
	channelID := item.Snippet.ChannelId
	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		log.Fatalf("%v", err)
	}

	video := Video{
		VideoID:      videoID,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		ChannelID:    channelID,
		PublishedAt:  publishedAt,
	}
	return video
}

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
func (v *Video) Insert() error {
	db := mydb.NewGormConnect()
	defer db.Close()

	err := v.SetDetailInfo()
	if err != nil {
		log.Printf("%v", err)
		return nil
	}

	db.Create(&v)
	log.Printf("Insert video: %v\n", v)

	return nil
}

// Update video
func (v *Video) Update() error {
	db := mydb.NewGormConnect()
	defer db.Close()

	err := v.SetDetailInfo()
	if err != nil {
		log.Printf("%v", err)
		return nil
	}

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

	return nil
}

// Delete video
func (v *Video) Delete() {
	db := mydb.NewGormConnect()
	defer db.Close()

	db.Delete(&v)

	log.Printf("Delete video: %v %v\n", v.VideoID, v.Title)
}

// SetDetailInfo ViewCount, CommentCount, CategoryID, CategoryName
func (v *Video) SetDetailInfo() error {
	service := NewYoutubeService()
	call := service.Videos.List("snippet,Statistics").
		Id(v.VideoID).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	} else if len(response.Items) == 0 {
		return youtubeError{
			content: "video",
			id:      v.VideoID,
			message: "This video has been deleted",
		}
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

	return nil
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
		log.Printf("%v", err)
		return nil
	}

	comments := []Comment{}
	for _, item := range response.Items {
		comments = append(comments, newComment(*item))
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

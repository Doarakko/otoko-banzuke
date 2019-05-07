package youtube

import (
	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	"google.golang.org/api/youtube/v3"
	"log"
	"time"
)

// Channel struct
type Channel struct {
	ChannelID       string    `gorm:"column:channel_id;primary_key"`
	Name            string    `gorm:"column:name;not null"`
	Description     string    `gorm:"column:description;not null"`
	ThumbnailURL    string    `gorm:"column:thumbnail_url;not null"`
	ViewCount       int64     `gorm:"column:view_count;not null"`
	VideoCount      int32     `gorm:"column:video_count;not null"`
	SubscriberCount int32     `gorm:"column:subscriber_count;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
	ExistChannel    bool      `gorm:"-"`
}

// Exists if channel exist return true
func (c *Channel) Exists() bool {
	db := mydb.NewGormConnect()
	defer db.Close()

	result := db.First(&c, "channel_id=?", c.ChannelID)

	return result.Error == nil
}

// Insert Channel
func (c *Channel) Insert() error {
	db := mydb.NewGormConnect()
	defer db.Close()

	c.SetDetailInfo()

	r := db.Create(&c)
	log.Printf("Insert channel: %v\n", r)

	return r.Error
}

// Update channel
func (c *Channel) Update() {
	db := mydb.NewGormConnect()
	defer db.Close()

	c.SetDetailInfo()

	db.Model(&c).Updates(Channel{
		Name:            c.Name,
		Description:     c.Description,
		ThumbnailURL:    c.ThumbnailURL,
		ViewCount:       c.ViewCount,
		SubscriberCount: c.SubscriberCount,
		VideoCount:      c.VideoCount,
	})
	log.Printf("Update channel: %v\n", c.ChannelID)
}

// Delete channel
func (c *Channel) Delete() {

}

func (c *Channel) selectVideos() []Video {
	db := mydb.NewGormConnect()
	defer db.Close()

	videos := []Video{}
	db.Find(&videos, "channel_id=?", c.ChannelID)

	return videos
}

// SetDetailInfo ViewCount, SubscriberCount, VideoCount
func (c *Channel) SetDetailInfo() {
	service := NewYoutubeService()
	call := service.Channels.List("snippet,statistics").
		Id(c.ChannelID).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	item := response.Items[0]

	c.Name = item.Snippet.Title
	c.Description = item.Snippet.Description
	c.ThumbnailURL = item.Snippet.Thumbnails.High.Url
	c.ViewCount = int64(item.Statistics.ViewCount)
	c.SubscriberCount = int32(item.Statistics.SubscriberCount)
	c.VideoCount = int32(item.Statistics.VideoCount)
}

// GetNewVideos assume to run once a day
func (c *Channel) GetNewVideos() []Video {
	// put 1 day period afer video published
	beginAt := time.Now().Add(-time.Duration(24*2) * time.Hour).Format(time.RFC3339)
	endAt := time.Now().Add(-time.Duration(24) * time.Hour).Format(time.RFC3339)

	service := NewYoutubeService()
	call := service.Search.List("id,snippet").
		Type("video").
		ChannelId(c.ChannelID).
		PublishedAfter(beginAt).
		PublishedBefore(endAt).
		Order("date").
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []Video{}
	for _, item := range response.Items {
		videos = append(videos, newVideo(*item))
	}
	log.Printf("Get %v new videos\n", len(videos))

	return videos
}

// GetAllVideos get all videos
func (c *Channel) GetAllVideos(pageToken string) []Video {
	service := NewYoutubeService()
	call := service.Search.List("id,snippet").
		Type("video").
		ChannelId(c.ChannelID).
		PageToken(pageToken).
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []Video{}
	for _, item := range response.Items {
		videos = append(videos, newVideo(*item))
	}

	pageToken = response.NextPageToken
	if pageToken != "" {
		videos = append(videos, c.GetAllVideos(pageToken)...)
	}
	log.Printf("Get %v videos\n", len(videos))

	return videos
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

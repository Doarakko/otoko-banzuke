package youtube

import (
	"log"
	"time"
)

// Channel fa
type Channel struct {
	ChannelID       string    `gorm:"column:channel_id"`
	Name            string    `gorm:"column:name"`
	Description     string    `gorm:"column:description"`
	ThumbnailURL    string    `gorm:"column:thumbnail_url"`
	PlaylistID      string    `gorm:"column:playlist_id"`
	ViewCount       int64     `gorm:"column:view_count"`
	VideoCount      int32     `gorm:"column:video_count"`
	SubscriberCount int32     `gorm:"column:subscriber_count"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (c *Channel) selectChannel() Channel {
	db := newGormConnect()
	defer db.Close()

	db.First(&c, "channel_id=?", c.ChannelID)

	return *c
}

// Insert hoge
func (c *Channel) Insert() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&c)
	log.Printf("Insert channel: %v\n", c)
}

func (c *Channel) update() {

}

func (c *Channel) delete() {

}

func (c *Channel) selectVideos() []Video {
	db := newGormConnect()
	defer db.Close()

	videos := []Video{}
	db.Find(&videos, "channel_id=?", c.ChannelID)

	return videos
}

func (c *Channel) deleteVideos() {

}

// SetDetailInfo hoge
func (c *Channel) SetDetailInfo() {
	service := newYoutubeService()
	call := service.Channels.List("snippet,contentDetails,statistics").
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
	c.PlaylistID = item.ContentDetails.RelatedPlaylists.Uploads
	c.ViewCount = int64(item.Statistics.ViewCount)
	c.SubscriberCount = int32(item.Statistics.SubscriberCount)
	c.VideoCount = int32(item.Statistics.VideoCount)
	// commentCount := item.Statistics.CommentCount
}

// TODO 日付指定
func (c *Channel) getNewVideos() []Video {
	service := newYoutubeService()
	call := service.Search.List("id").
		Type("video").
		ChannelId(c.ChannelID).
		Order("date").
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []Video{}
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

		video := Video{
			VideoID:      videoID,
			Title:        title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
			ChannelID:    channelID,
			PublishedAt:  publishedAt,
		}
		videos = append(videos, video)
	}
	log.Printf("Get %v videos\n", len(videos))

	return videos
}

func getAllVideos(playlistID string, pageToken string) []Video {
	service := newYoutubeService()
	call := service.PlaylistItems.List("id,snippet,contentDetails").
		PlaylistId(playlistID).
		PageToken(pageToken).
		MaxResults(50)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	videos := []Video{}
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

		video := Video{
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
	log.Printf("Get %v videos\n", len(videos))

	return videos
}

func getHighRatedVideos(playlistID string, pageToken string) {

}

// SearchChannels hoge
func SearchChannels(q string) []Channel {
	service := newYoutubeService()
	call := service.Search.List("id,snippet").
		Type("channel").
		Q(q).
		Order("relevance").
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	channels := []Channel{}
	for _, item := range response.Items {
		channelID := item.Id.ChannelId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url

		channel := Channel{
			ChannelID:    channelID,
			Name:         title,
			Description:  description,
			ThumbnailURL: thumbnailURL,
		}
		channels = append(channels, channel)
	}
	log.Printf("Get %v channels\n", len(channels))

	return channels
}

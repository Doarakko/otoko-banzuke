package youtube

import (
	"log"
	"time"
)

// Channel fa
type Channel struct {
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

func (c *Channel) selectChannel() Channel {
	db := newGormConnect()
	defer db.Close()

	db.First(&c, "channel_id=?", c.ChannelID)

	return *c
}

func (c *Channel) insertChannel() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&c)
	log.Printf("Insert channel: %v\n", c)
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

func (c *Channel) getChannelInfo() Channel {
	service := newYoutubeService()
	call := service.Channels.List("snippet,contentDetails,statistics").
		Id(c.ChannelID).
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

	*c = Channel{
		Name:            name,
		Description:     description,
		ThumbnailURL:    thumbnailURL,
		PlaylistID:      playlistID,
		ViewCount:       viewCount,
		VideoCount:      videoCount,
		SubscriberCount: subscriberCount,
	}
	return *c
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

// SearchChannel hoge
func SearchChannel(q string) []Channel {
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
			Description:  description,
			ThumbnailURL: thumbnailURL,
			Name:         title,
			ChannelID:    channelID,
		}
		channels = append(channels, channel)
	}
	log.Printf("Get %v channel\n", len(channels))

	return channels
}

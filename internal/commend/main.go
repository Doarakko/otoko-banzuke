package commend

import (
	"log"

	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SearchChannels hoge
func SearchChannels(q string) []myyoutube.Channel {
	service := myyoutube.NewYoutubeService()
	call := service.Search.List("id,snippet").
		Type("channel").
		Q(q).
		Order("relevance").
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	channels := []myyoutube.Channel{}
	for _, item := range response.Items {
		channelID := item.Id.ChannelId
		title := item.Snippet.Title
		description := item.Snippet.Description
		thumbnailURL := item.Snippet.Thumbnails.High.Url

		channel := myyoutube.Channel{
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

// InsertChannel get channel detail information and insert channel
func InsertChannel(channelID string) {
	channel := myyoutube.Channel{
		ChannelID: channelID,
	}
	channel.SetDetailInfo()
	channel.Insert()
}

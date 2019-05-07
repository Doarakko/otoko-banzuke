package new

import (
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SelectNewComments get new comments
func SelectNewComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}

	// put 1 day period afer video published where get new video
	preAt := time.Now().Add(-time.Duration(24*7) * time.Hour)

	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, comments.channel_id, videos.thumbnail_url, channels.name, RANK() OVER (ORDER BY comments.like_count DESC)").
		Where("videos.published_at >= ?", preAt).
		Joins("JOIN videos ON videos.video_id = comments.video_id").
		Joins("JOIN channels ON channels.channel_id = comments.channel_id").
		Order("rank").
		Find(&comments)

	return comments
}

package rank

import (
	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SelectRankComments get comments and rank based on like count
func SelectRankComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}
	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, comments.channel_id, videos.thumbnail_url, channels.name, RANK() OVER (ORDER BY comments.like_count DESC)").
		Joins("JOIN videos ON videos.video_id = comments.video_id").
		Joins("JOIN channels ON channels.channel_id = comments.channel_id").
		Order("rank").
		Limit(48).
		Find(&comments)

	return comments
}

package search

import (
	"strings"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SearchComment search comment related with parameter from database
func SearchComment(q string) []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}
	q = "%" + strings.ToLower(q) + "%"
	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, videos.thumbnail_url, channels.name, comments.channel_id").
		Joins("JOIN videos ON videos.video_id = comments.video_id").
		Joins("JOIN channels ON channels.channel_id = comments.channel_id").
		Where("lower(comments.text_display) like ? or lower(channels.name) like ? or lower(videos.title) like ?", q, q, q).
		Order("comments.like_count desc").
		Find(&comments)

	return comments
}

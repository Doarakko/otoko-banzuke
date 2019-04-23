package search

import (
	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SearchOtoko hoge
func SearchOtoko(q string) []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}
	q = "%" + q + "%"
	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, videos.thumbnail_url").
		Where("text_display like ?", q).Joins("JOIN videos ON videos.video_id = comments.video_id").
		Order("comments.like_count desc").
		Find(&comments)

	return comments
}

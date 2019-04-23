package banzuke

import (
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SelectRankComments get comments
func SelectRankComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}
	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, videos.thumbnail_url, RANK() OVER (ORDER BY comments.like_count DESC)").
		Order("rank").
		Joins("JOIN videos ON videos.video_id = comments.video_id").
		Find(&comments)

	return comments
}

// SelectTodayComments get today comments
func SelectTodayComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}

	preAt := time.Now().Add(-time.Duration(24 * time.Hour))
	db.Where("created_at >= ?", preAt).Order("like_count desc").Find(&comments)

	return comments
}

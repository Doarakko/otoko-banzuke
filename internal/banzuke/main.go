package banzuke

import (
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

// SelectRankComments get comments and rank based on like count
func SelectRankComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}
	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, videos.thumbnail_url, RANK() OVER (ORDER BY comments.like_count DESC)").
		Joins("JOIN videos ON videos.video_id = comments.video_id").
		Order("rank").
		Limit(50).
		Find(&comments)

	return comments
}

// SelectTodayComments get today comments
func SelectTodayComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}

	// put 1 day period afer video published where get new video
	preAt := time.Now().Add(-time.Duration(24*2) * time.Hour)

	db.Select("comments.text_display, comments.like_count, comments.video_id, comments.author_name, comments.author_url, videos.thumbnail_url").
		Where("videos.published_at >= ?", preAt).
		Joins("JOIN videos ON videos.video_id = comments.video_id").
		Order("comments.like_count DESC").
		Limit(5).
		Find(&comments)

	return comments
}

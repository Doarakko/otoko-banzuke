package base

import (
	mydb "github.com/Doarakko/otoko-banzuke/database"
)

// GetTotalComment hoge
func GetTotalComment() int {
	db := mydb.NewGormConnect()
	defer db.Close()

	var count int
	db.Table("comments").
		Select("count(distinct(comments.comment_id))").
		Count(&count)

	return count
}

// GetTotalAuthor hoge
func GetTotalAuthor() int {
	db := mydb.NewGormConnect()
	defer db.Close()

	var count int
	db.Table("comments").
		Select("count(distinct(comments.author_url))").
		Count(&count)

	return count
}

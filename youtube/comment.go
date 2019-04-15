package main

import (
	"log"
	"time"
)

type comment struct {
	CommentID   string
	TextDisplay string
	AuthorID    string
	AuthorName  string
	AuthorURL   string
	ChannelID   string
	VideoID     string
	ParentID    string
	LikeCount   int32
	ReplyCount  int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c *comment) insertComment() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&c)
	log.Printf("Insert comment: %v\n", c)
}

func (c *comment) checkOtoko() {

}

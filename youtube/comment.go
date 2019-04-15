package main

import (
	"log"
	"regexp"
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

var res = []*regexp.Regexp{
	regexp.MustCompile("^.+男.{0,1}$"),
	regexp.MustCompile("^.+漢.{0,1}$"),
	regexp.MustCompile("^.+おとこ.{0,1}$"),
	regexp.MustCompile("^.+オトコ.{0,1}$"),
	regexp.MustCompile("^.+女.{0,1}$"),
	regexp.MustCompile("^.+おんな.{0,1}$"),
	regexp.MustCompile("^.+オンナ.{0,1}$"),
}

func (c *comment) insertComment() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&c)
	log.Printf("Insert comment: %v\n", c)
}

func (c *comment) checkOtoko() bool {
	for _, re := range res {
		if re.MatchString(c.TextDisplay) {
			return true
		}
	}
	return false
}

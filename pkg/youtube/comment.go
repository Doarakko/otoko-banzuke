package youtube

import (
	"log"
	"regexp"
	"time"
)

// Comment gaerg
type Comment struct {
	CommentID   string    `gorm:"column:comment_id"`
	TextDisplay string    `gorm:"column:text_display"`
	AuthorName  string    `gorm:"column:author_name"`
	AuthorURL   string    `gorm:"column:author_url"`
	LikeCount   int32     `gorm:"column:like_count"`
	ReplyCount  int32     `gorm:"column:reply_count"`
	ChannelID   string    `gorm:"column:channel_id"`
	VideoID     string    `gorm:"column:video_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
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

func (c *Comment) insert() {
	db := newGormConnect()
	defer db.Close()

	db.Create(&c)
	log.Printf("Insert comment: %v\n", c)
}

func (c *Comment) update() {

}

func (c *Comment) delete() {

}

func (c *Comment) getCommentInfo() Comment {
	return *c
}

func (c *Comment) checkOtoko() bool {
	for _, re := range res {
		if re.MatchString(c.TextDisplay) {
			return true
		}
	}
	return false
}

// SelectAllComments get comments
func SelectAllComments() []Comment {
	db := newGormConnect()
	defer db.Close()

	comments := []Comment{}
	db.Find(&comments)

	return comments
}

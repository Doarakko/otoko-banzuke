package youtube

import (
	"log"
	"regexp"
	"time"
)

// Comment gaerg
type Comment struct {
	CommentID    string
	TextDisplay  string
	AuthorID     string
	AuthorName   string
	AuthorURL    string
	ChannelID    string
	VideoID      string
	ThumbnailURL string
	ParentID     string
	LikeCount    int32
	ReplyCount   int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

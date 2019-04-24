package youtube

import (
	"log"
	"regexp"
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
)

// Comment struct
type Comment struct {
	CommentID    string    `gorm:"column:comment_id"`
	TextDisplay  string    `gorm:"column:text_display"`
	AuthorName   string    `gorm:"column:author_name"`
	AuthorURL    string    `gorm:"column:author_url"`
	LikeCount    int32     `gorm:"column:like_count"`
	ReplyCount   int32     `gorm:"column:reply_count"`
	ChannelID    string    `gorm:"column:channel_id"`
	VideoID      string    `gorm:"column:video_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	Rank         int       `gorm:"-"`
	ThumbnailURL string    `gorm:"-"`
}

// Insert comment
func (c *Comment) Insert() {
	db := mydb.NewGormConnect()
	defer db.Close()

	db.Create(&c)
	log.Printf("Insert comment: %v\n", c)
}

// Update comment
func (c *Comment) Update() {

}

// Delete comment
func (c *Comment) Delete() {

}

var re = regexp.MustCompile("^.+(男|漢|おとこ|オトコ|女|おんな|オンナ).{0,1}$")

// CheckOtoko if otoko comment return true
func (c *Comment) CheckOtoko() bool {
	return re.MatchString(c.TextDisplay)
}

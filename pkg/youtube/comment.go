package youtube

import (
	"log"
	"regexp"
	"time"
)

// Comment gaerg
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

var re = regexp.MustCompile("^.+(男|漢|おとこ|オトコ|女|おんな|オンナ).{0,1}$")

// Insert comment
func (c *Comment) Insert() {
	db := newGormConnect()
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

// CheckOtoko oge
func (c *Comment) CheckOtoko() bool {
	if re.MatchString(c.TextDisplay) {
		return true
	}
	return false
}

// SelectRankComments get comments
func SelectRankComments() []Comment {
	db := newGormConnect()
	defer db.Close()

	comments := []Comment{}
	db.Select("*, RANK() OVER (ORDER BY comments.like_count DESC)").Order("rank").Joins("JOIN videos ON videos.video_id = comments.video_id").Find(&comments)

	return comments
}

// SelectTodayComments get today comments
func SelectTodayComments() []Comment {
	db := newGormConnect()
	defer db.Close()

	comments := []Comment{}

	preAt := time.Now().Add(-time.Duration(24 * time.Hour))
	db.Where("created_at >= ?", preAt).Order("like_count desc").Find(&comments)

	return comments
}

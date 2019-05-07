package youtube

import (
	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	"google.golang.org/api/youtube/v3"
	"log"
	"regexp"
	"time"
)

// Comment struct
type Comment struct {
	CommentID    string    `gorm:"column:comment_id;primary_key"`
	TextDisplay  string    `gorm:"column:text_display;not null"`
	AuthorName   string    `gorm:"column:author_name;not null"`
	AuthorURL    string    `gorm:"column:author_url;not null"`
	LikeCount    int32     `gorm:"column:like_count;not null"`
	ReplyCount   int32     `gorm:"column:reply_count;not null"`
	ChannelID    string    `gorm:"column:channel_id;not null;index"`
	VideoID      string    `gorm:"column:video_id;not null;index"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
	Rank         int       `gorm:"-"`
	Name         string    `gorm:"-"`
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
	db := mydb.NewGormConnect()
	defer db.Close()

	c.SetDetailInfo()

	db.Model(&c).Updates(Comment{
		LikeCount:   c.LikeCount,
		AuthorName:  c.AuthorName,
		TextDisplay: c.TextDisplay,
	})
	log.Printf("Update comment: %v\n", c.CommentID)
}

// Delete comment
func (c *Comment) Delete() {

}

// SetDetailInfo ViewCount, CommentCount, CategoryID, CategoryName
func (c *Comment) SetDetailInfo() {
	service := NewYoutubeService()
	call := service.Comments.List("id,snippet").
		Id(c.CommentID).
		TextFormat("plainText").
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	item := response.Items[0]

	c.AuthorName = item.Snippet.AuthorDisplayName
	c.TextDisplay = item.Snippet.TextDisplay
	c.LikeCount = int32(item.Snippet.LikeCount)
}

var re = regexp.MustCompile("^.+(男|漢|おとこ|オトコ|女|おんな|オンナ).{0,1}$")

// CheckComment if otoko comment return true
func (c *Comment) CheckComment() bool {
	return re.MatchString(c.TextDisplay)
}

// Reply comment
func (c *Comment) Reply() {
	reply := &youtube.Comment{
		Snippet: &youtube.CommentSnippet{
			ParentId:     c.CommentID,
			TextOriginal: "情報提供ありがとうございます。\n漢番付に登録しました。\n【漢番付】https://otoko-banzuke.herokuapp.com/",
		},
	}

	service := NewYoutubeService()
	call := service.Comments.Insert("id", reply)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Reply to %v, from %v\n", c.CommentID, response.Id)
}

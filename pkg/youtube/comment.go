package youtube

import (
	"log"
	"regexp"
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	"google.golang.org/api/youtube/v3"
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
func (c *Comment) Update() error {
	db := mydb.NewGormConnect()
	defer db.Close()

	err := c.SetDetailInfo()
	if err != nil {
		log.Printf("%v", err)
		return nil
	}

	db.Model(&c).Updates(Comment{
		LikeCount:   c.LikeCount,
		AuthorName:  c.AuthorName,
		TextDisplay: c.TextDisplay,
	})
	log.Printf("Update comment: %v %v %v\n", c.CommentID, c.TextDisplay, c.LikeCount)
	return nil
}

// Delete comment
func (c *Comment) Delete() {
	db := mydb.NewGormConnect()
	defer db.Close()

	db.Delete(&c)

	log.Printf("Delete comment: %v %v\n", c.CommentID, c.TextDisplay)
}

// SetDetailInfo AuthorName, TextDisplay, LikeCount
func (c *Comment) SetDetailInfo() error {
	service := NewYoutubeService()
	call := service.Comments.List("snippet").
		Id(c.CommentID).
		TextFormat("plainText").
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	} else if len(response.Items) == 0 {
		return youtubeError{
			content: "comment",
			id:      c.CommentID,
			message: "This comment has been deleted",
		}
	}
	item := response.Items[0]

	c.AuthorName = item.Snippet.AuthorDisplayName
	c.TextDisplay = item.Snippet.TextDisplay
	c.LikeCount = int32(item.Snippet.LikeCount)

	return nil
}

var re = regexp.MustCompile("^.+(漢|漢達|男|男達|男性|おとこ|オトコ|女|女達|女性|おんな|オンナ)(。|\\.|~|〜|!|！|\\*|＊|w|W|♂|♀){0,1}$")

// CheckComment if otoko comment return true
func (c *Comment) CheckComment() bool {
	return re.MatchString(c.TextDisplay) && c.LikeCount >= 5
}

func newComment(item youtube.CommentThread, channelID string) Comment {
	commentID := item.Snippet.TopLevelComment.Id
	videoID := item.Snippet.VideoId
	authorName := item.Snippet.TopLevelComment.Snippet.AuthorDisplayName
	authorURL := item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl
	textDisplay := item.Snippet.TopLevelComment.Snippet.TextDisplay
	likeCount := int32(item.Snippet.TopLevelComment.Snippet.LikeCount)
	replyCount := int32(item.Snippet.TotalReplyCount)

	comment := Comment{
		CommentID:   commentID,
		VideoID:     videoID,
		ChannelID:   channelID,
		TextDisplay: textDisplay,
		AuthorName:  authorName,
		AuthorURL:   authorURL,
		LikeCount:   likeCount,
		ReplyCount:  replyCount,
	}
	return comment
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

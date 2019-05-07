package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	//"github.com/joho/godotenv"

	rank "github.com/Doarakko/otoko-banzuke/internal/rank"
	base "github.com/Doarakko/otoko-banzuke/internal/base"
	commend "github.com/Doarakko/otoko-banzuke/internal/commend"
	search "github.com/Doarakko/otoko-banzuke/internal/search"
	new "github.com/Doarakko/otoko-banzuke/internal/new"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
)

func main() {
	// err := godotenv.Load("./.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	router := gin.Default()
	router.Static("web/static", "./web/static")
	router.LoadHTMLFiles(
		"web/template/index.tmpl", 
		"web/template/new/index.tmpl",
		"web/template/commend/index.tmpl", 
		"web/template/search/index.tmpl",
	)

	totalComment := base.GetTotalComment()
	totalAuthor := base.GetTotalAuthor()
	rankComments := rank.SelectRankComments()
	newComments := new.SelectNewComments()

	// 番付
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"totalComment":  totalComment,
			"totalAuthor":   totalAuthor,
			"rankComments":  rankComments,
		})
	})

	// 今週の漢
	router.GET("/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new/index.tmpl", gin.H{
			"totalComment":  totalComment,
			"totalAuthor":   totalAuthor,
			"newComments": newComments,
		})
	})

	// 漢を推薦する
	// initial display
	channels := []myyoutube.Channel{}
	router.GET("/commend", func(c *gin.Context) {
		c.HTML(http.StatusOK, "commend/index.tmpl", gin.H{
			"totalComment": totalComment,
			"totalAuthor":  totalAuthor,
			"channels":     channels,
		})
	})

	router.POST("/commend", func(c *gin.Context) {
		query := c.PostForm("query")
		channelID := c.PostForm("channel_id")

		if query != "" {
			channels = commend.SearchChannels(query)
			for i := range channels {
				channels[i].SetDetailInfo()
			}
		} else if channelID != "" {
			commend.InsertChannel(channelID)
		}

		c.Redirect(302, "/commend")
	})

	// 漢を探す
	// initial display
	comments := []myyoutube.Comment{}
	commentCount := len(comments)
	router.GET("/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
			"totalComment": totalComment,
			"totalAuthor":  totalAuthor,
			"comments":     comments,
			"commentCount": commentCount,
		})
	})

	router.POST("/search", func(c *gin.Context) {
		query := c.PostForm("query")
		comments = search.SearchComment(query)
		commentCount = len(comments)

		c.Redirect(302, "/search")
	})

	router.Run()
}

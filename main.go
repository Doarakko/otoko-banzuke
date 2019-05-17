package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	base "github.com/Doarakko/otoko-banzuke/internal/base"
	commend "github.com/Doarakko/otoko-banzuke/internal/commend"
	new "github.com/Doarakko/otoko-banzuke/internal/new"
	rank "github.com/Doarakko/otoko-banzuke/internal/rank"
	search "github.com/Doarakko/otoko-banzuke/internal/search"
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

	// 番付
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"totalComment": base.GetTotalComment(),
			"totalAuthor":  base.GetTotalAuthor(),
			"rankComments": rank.SelectRankComments(),
		})
	})

	// 今週の漢
	router.GET("/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new/index.tmpl", gin.H{
			"totalComment": base.GetTotalComment(),
			"totalAuthor":  base.GetTotalAuthor(),
			"newComments":  new.SelectNewComments(),
		})
	})

	// 漢を推薦する
	router.GET("/commend", func(c *gin.Context) {
		c.HTML(http.StatusOK, "commend/index.tmpl", gin.H{
			"totalComment": base.GetTotalComment(),
			"totalAuthor":  base.GetTotalAuthor(),
		})
	})

	router.POST("/commend", func(c *gin.Context) {
		query := c.PostForm("query")
		channelID := c.PostForm("channel_id")

		if channelID != "" {
			commend.InsertChannel(channelID)
			c.Redirect(302, "/commend")
		} else if query != "" {
			c.HTML(http.StatusOK, "commend/index.tmpl", gin.H{
				"totalComment": base.GetTotalComment(),
				"totalAuthor":  base.GetTotalAuthor(),
				"channels":     commend.SearchChannels(query),
			})
		}
	})

	// 漢を探す
	router.GET("/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
			"totalComment": base.GetTotalComment(),
			"totalAuthor":  base.GetTotalAuthor(),
			"commentCount": -1,
		})
	})

	router.POST("/search", func(c *gin.Context) {
		query := c.PostForm("query")
		comments := search.SearchComment(query)

		c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
			"totalComment": base.GetTotalComment(),
			"totalAuthor":  base.GetTotalAuthor(),
			"comments":     comments,
			"commentCount": len(comments),
		})
	})

	router.Run()
}

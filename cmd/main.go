package main

import (
	"log"
	"net/http"

	myyoutube "../pkg/youtube"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	comments := myyoutube.SelectAllComments()

	router := gin.Default()
	router.Static("/web/static", "../web/static")
	router.LoadHTMLFiles("../web/template/index.tmpl", "../web/template/commend/index.tmpl")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"comments": comments,
		})
	})

	channels := []myyoutube.Channel{}
	router.GET("/commend", func(c *gin.Context) {
		c.HTML(http.StatusOK, "commend/index.tmpl", gin.H{
			"channels": channels,
		})
	})

	router.POST("/commend", func(c *gin.Context) {
		query := c.PostForm("query")
		channelID := c.PostForm("channel_id")

		if query != "" {
			channels = myyoutube.SearchChannels(query)
			for i := range channels {
				channels[i].SetDetailInfo()
			}
		} else {
			channel := myyoutube.Channel{
				ChannelID: channelID,
			}
			channel.SetDetailInfo()
			channel.Insert()
		}
		c.Redirect(302, "/commend")
	})

	router.Run(":8080")
}

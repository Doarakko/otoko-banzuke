package main

import (
	"log"
	"net/http"

	"./youtube"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	comments := youtube.SelectAllComments()

	router := gin.Default()
	router.Static("/public", "./public")
	router.LoadHTMLFiles("templates/index.tmpl", "templates/form/index.tmpl")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"comments": comments,
		})
	})

	channels := []youtube.Channel{}
	router.GET("/form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form/index.tmpl", gin.H{
			"channels": channels,
		})
	})

	router.POST("/form", func(c *gin.Context) {
		query := c.PostForm("query")
		channels = youtube.SearchChannel(query)
		c.Redirect(302, "/form")
	})

	router.Run(":8080")
}

package main

import (
	"log"
	"net/http"

	"./youtube"
	"github.com/gin-gonic/gin"
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
	router.LoadHTMLGlob("templates/index.tmpl")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "Main website",
			"comments": comments,
		})
	})

	//router.POST("/index", SignupedRoute)

	router.Run(":8080")
}

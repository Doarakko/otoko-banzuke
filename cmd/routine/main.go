package main

import (
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func updateAllComments() {

}

func main() {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	for _, channel := range myyoutube.SelectAllChannels() {
		for _, video := range channel.GetNewVideos() {
			for _, comment := range video.GetComments() {
				if comment.CheckOtoko() {
					comment.Insert()
				}
			}
		}
	}
}

package main

import (
	//"log"
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	//"github.com/joho/godotenv"
)

func updateAllChannels() {
	for _, channel := range selectAllChannels() {
		channel.Update()
	}
}

func updateAllVideos() {
	for _, video := range selectAllVideos() {
		video.Update()
	}
}

func updateAllComments() {
	for _, comment := range selectAllComments() {
		comment.Update()
	}
}

func selectAllChannels() []myyoutube.Channel {
	db := mydb.NewGormConnect()
	defer db.Close()

	channels := []myyoutube.Channel{}
	db.Find(&channels)

	return channels
}

func selectAllVideos() []myyoutube.Video {
	db := mydb.NewGormConnect()
	defer db.Close()

	videos := []myyoutube.Video{}
	db.Find(&videos)

	return videos
}

func selectAllComments() []myyoutube.Comment {
	db := mydb.NewGormConnect()
	defer db.Close()

	comments := []myyoutube.Comment{}
	db.Find(&comments)

	return comments
}

func selectNewChannels() []myyoutube.Channel {
	db := mydb.NewGormConnect()
	defer db.Close()

	channels := []myyoutube.Channel{}

	preAt := time.Now().Add(-time.Duration(24 * time.Hour))
	db.Where("created_at >= ?", preAt).Find(&channels)

	return channels
}

func searchNewOtoko() {
	for _, channel := range selectAllChannels() {
		for _, video := range channel.GetNewVideos() {
			if video.Exists() {
				continue
			}

			createdFlag := false
			for _, comment := range video.GetComments() {
				if comment.CheckOtoko() {
					if !createdFlag {
						video.Insert()
					}
					comment.Insert()
					createdFlag = true
				}
			}
		}
	}
}

func searchAllOtoko() {
	for _, channel := range selectNewChannels() {
		for _, video := range channel.GetAllVideos("") {
			if video.Exists() {
				continue
			}

			createdFlag := false
			for _, comment := range video.GetComments() {
				if comment.CheckOtoko() {
					if !createdFlag {
						video.Insert()
					}
					comment.Insert()
					createdFlag = true
				}
			}
		}
	}
}

func main() {
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	searchAllOtoko()
	searchNewOtoko()
	//updateAllChannels()
	//updateAllVideos()
	//updateAllComments()
}

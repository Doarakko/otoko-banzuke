package main

import (
	"log"
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func updateAllChannels() {

}

func updateAllVideos() {

}

func updateAllComments() {

}

func selectAllChannels() []myyoutube.Channel {
	db := mydb.NewGormConnect()
	defer db.Close()

	channels := []myyoutube.Channel{}
	db.Find(&channels)

	return channels
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
			createdFlag := false
			log.Printf("========================start video===========%v=============\n", video.VideoID)
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
	searchAllOtoko()
	searchNewOtoko()
}

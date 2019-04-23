package main

import (
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/pkg/database"
	myyoutube "github.com/Doarakko/otoko-banzuke/pkg/youtube"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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

func updateAllComments() {

}

func updateAllChannels() {

}

func searchNewOtoko() {
	for _, channel := range selectAllChannels() {
		for _, video := range channel.GetNewVideos() {
			for _, comment := range video.GetComments() {
				if comment.CheckOtoko() {
					video.Insert()
					comment.Insert()
				}
			}
		}
	}
}

func searchAllOtoko() {
	for _, channel := range selectAllChannels() {
		for _, video := range channel.GetAllVideos("") {
			for _, comment := range video.GetComments() {
				if comment.CheckOtoko() {
					video.Insert()
					comment.Insert()
				}
			}
		}
	}
}

func main() {
	updateAllChannels()

}

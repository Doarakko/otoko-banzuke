package main

import (
	"log"
	"time"

	mydb "github.com/Doarakko/otoko-banzuke/database"
	"github.com/Doarakko/otoko-banzuke/slack"
	myyoutube "github.com/Doarakko/otoko-banzuke/youtube"

	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	preAt := time.Now().Add(-time.Duration(24) * time.Hour)
	log.Printf("%v\n", preAt)
	db.Where("created_at >= ?", preAt).Find(&channels)

	return channels
}

// Search and insert comments from new videos
func searchNewComments() {
	for _, channel := range selectAllChannels() {
		for _, video := range channel.GetNewVideos() {
			searchVideoComments(video)
		}
	}
}

// Search and insert comments from all videos
func searchAllComments() {
	for _, channel := range selectNewChannels() {
		for _, video := range channel.GetAllVideos("") {
			searchVideoComments(video)
		}
	}
}

// Search and insert comments from one video
func searchVideoComments(video myyoutube.Video) bool {
	if video.Exists() {
		log.Printf("skip video: %v", video.VideoID)
		return false
	}

	createdFlag := false
	for _, comment := range video.GetComments() {
		if comment.CheckComment() {
			if !createdFlag {
				video.Insert()
			}
			comment.Insert()
			createdFlag = true
		}
	}
	return true
}

func deleteComments() {
	for _, comment := range selectAllComments() {
		if !comment.CheckComment() {
			comment.Delete()
		}
	}
}

func main() {
	// err := godotenv.Load("./.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	//searchAllComments()
	searchNewComments()
	//updateAllChannels()
	//updateAllVideos()
	//updateAllComments()
	// channel := myyoutube.Channel{
	// 	ChannelID: "UCMJiPpN_09F0aWpQrgbc_qg",
	// }
	// for _, video := range channel.GetAllVideos("") {
	// 	searchVideoComments(video)
	// }
	// updateAllComments()
	slack.Post("Complete daily routine")
}

package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func newGormConnect() *gorm.DB {
	//db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := gorm.Open("postgres", "postgres://uojcpxxbqfdtwe:32aaebe5b808ffc63f7753a98f16000479469eb57e7a85eed218feb2e0436463@ec2-23-23-173-30.compute-1.amazonaws.com:5432/d16518rirmmbsl")

	if err != nil {
		panic(err.Error())
	}
	return db
}

func newYoutubeService() *youtube.Service {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyD3BcAAS4dmpvUbVrT8K3U49d6PYxR8CTk"},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	return service
}

func main() {
	channel := channel{}
	channel.ChannelID = "UCZFxcWJS1_iVIFETARRRHZQ"
	//insertChannel(getChannelInfo(channelID))
	channel = channel.selectChannel()
	print(channel.PlaylistID)
	//insertVideos(getAllVideos(channel.PlaylistID, ""))
}

package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Payload slack
type payload struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
	Channel  string `json:"channel"`
}

// Post to slack
func Post(message string) {
	payload := payload{
		Text:     fmt.Sprintf("%s", message),
		Username: "漢番付",
		IconURL:  "https://otoko-banzuke.herokuapp.com/web/static/favicon.ico",
	}
	jsonParams, _ := json.Marshal(payload)
	_, err := http.PostForm(
		os.Getenv("SLACK_WEBHOOK_URL"),
		url.Values{"payload": {string(jsonParams)}},
	)
	if err != nil {
		log.Printf("%v", err)
	}
}

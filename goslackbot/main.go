package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {

	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	channel_id := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: "Super Bot Message",
		Text:    "Golang COnnection Test",
		Color:   "4af030",
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}

	_, timestamp, err := client.PostMessage(channel_id, slack.MsgOptionAttachments(attachment))

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent at %s", timestamp)

}

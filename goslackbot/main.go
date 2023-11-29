package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {

	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	app_token := os.Getenv("SLACK_APP_TOKEN")
	//channel_id := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(app_token))

	socket_client := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func(ctx context.Context, client *slack.Client, socket_client *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener.")
				return
			case event := <-socket_client.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					events_api, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast event to the events API: %v\n", event)
						continue
					}

					socket_client.Ack(*event.Request)
					log.Println(events_api)
				}
			}
		}

	}(ctx, client, socket_client)

	socket_client.Run()
}

// Slack上で動くbotくん
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type envConfig struct {
	// BotToken is bot user token to access to slack API.
	BotToken  string `envconfig:"BOT_TOKEN" required:"true"`
	RakutenID string `envconfig:"RAKUTEN_ID" required:"true"`
	BotName   string `envconfig:"BOT_NAME" required:"true"`
	BotID     string `envconfig:"BOT_ID" required:"true"`
}

var env envConfig

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(
		env.BotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	rtm := client.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			msg := handleMsg(ev.Text)
			if msg != "" {
				rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
			}
		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
			return 1

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return 1

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
	return 0
}

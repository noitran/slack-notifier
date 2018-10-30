package main

import (
	"encoding/json"
	"fmt"
	"github.com/int128/slack"
	"os"
	"time"
)

const (
	EnvWebhookUrl = "SLACK_WEBHOOK"
	EnvIcon = "SLACK_ICON"
	EnvChannel = "SLACK_CHANNEL"
	EnvUsername = "SLACK_USERNAME"
	EnvMessage = "SLACK_MESSAGE"
	EnvBody = "SLACK_ATTACHMENT"
)

func main() {
	webhookUrl := os.Getenv(EnvWebhookUrl)

	if webhookUrl == "" {
		stdErr("Webhook url is required")
	}

	attachmentStruct := slack.Attachment{
		Timestamp: time.Now().Unix(),
		MrkdwnIn:  []string{"text"},
	}

	body := []byte(os.Getenv(EnvBody))

 	err := json.Unmarshal(body, &attachmentStruct)
	if  err != nil {
		stdErr("Failed to parse attachment")
	}

	var msg = &slack.Message{
		Username: os.Getenv(EnvUsername),
		IconEmoji: os.Getenv(EnvIcon),
		Channel: os.Getenv(EnvChannel),
		Text: os.Getenv(EnvMessage),
		Attachments: []slack.Attachment{
			attachmentStruct,
		},
	}

	if err := slack.Send(webhookUrl, msg); err != nil {
		stdErr(fmt.Sprintf("Could not send the message to Slack: %s", err))
	}
}

func stdErr(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

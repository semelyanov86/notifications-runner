package service

import (
	"Notifications/internal/config"
	"Notifications/pkg/e"
	"fmt"
	"github.com/slack-go/slack"
)

type Slack struct {
	Config config.Config
}

func NewSlack(config config.Config) Slack {
	return Slack{
		Config: config,
	}
}

func (s Slack) Send(msg string) (string, string, error) {
	api := slack.New(s.Config.Slack.Token)

	channelId, timestamp, err := api.PostMessage(
		s.Config.Slack.ChannelId,
		slack.MsgOptionText(msg, false),
	)

	if err != nil {
		return "", "", e.Wrap("can not send slack message", err)
	}

	fmt.Printf("Message sent successfully to channel %s at %s", channelId, timestamp)
	return timestamp, channelId, nil
}

package service

import (
	"Notifications/internal/config"
	"Notifications/internal/domain"
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

func (s Slack) Send(entity domain.Entity) (string, string, error) {
	token := entity.KeyId
	if token == "" {
		token = s.Config.Slack.Token
	}
	channelId := s.Config.Slack.ChannelId
	if entity.ChatId.Valid {
		channelId = entity.ChatId.String
	}
	api := slack.New(token)

	channelId, timestamp, err := api.PostMessage(
		channelId,
		slack.MsgOptionText(entity.Description, false),
	)

	if err != nil {
		return "", "", e.Wrap("can not send slack message", err)
	}

	fmt.Printf("Message sent successfully to channel %s at %s", channelId, timestamp)
	fmt.Println()
	return timestamp, channelId, nil
}

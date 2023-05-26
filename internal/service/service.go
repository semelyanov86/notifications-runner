package service

import "github.com/semelyanov86/notifications-runner/internal/config"

type Services struct {
	Slack Slack
}

func NewServices(config config.Config) *Services {
	return &Services{Slack: NewSlack(config)}
}

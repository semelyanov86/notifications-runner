package service

import "Notifications/internal/config"

type Services struct {
	Slack Slack
}

func NewServices(config config.Config) *Services {
	return &Services{Slack: NewSlack(config)}
}

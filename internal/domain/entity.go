package domain

import "time"

type Entity struct {
	Crmid          int
	Label          string
	Description    string
	ChatId         string
	NotifyDateTime time.Time
	NotifyStatus   string
	NotifyType     string
}

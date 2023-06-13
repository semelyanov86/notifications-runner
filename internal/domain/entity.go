package domain

import (
	"database/sql"
	"time"
)

type Entity struct {
	Crmid          int
	Label          string
	Description    string
	ChatId         sql.NullString
	NotifyDateTime time.Time
	NotifyStatus   string
	NotifyType     string
	KeyId          string
}

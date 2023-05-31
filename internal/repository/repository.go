package repository

import (
	"Notifications/internal/config"
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Repositories struct {
	Entities Entity
}

func NewRepositories(db *sql.DB, config config.Config) *Repositories {
	return &Repositories{Entities: NewEntityRepository(db, config)}
}

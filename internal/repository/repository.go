package repository

import (
	"database/sql"
	"errors"
	"github.com/semelyanov86/notifications-runner/internal/config"
)

var ErrRecordNotFound = errors.New("record not found")

type Repositories struct {
	Entities Entity
}

func NewRepositories(db *sql.DB, config config.Config) *Repositories {
	return &Repositories{Entities: NewEntityRepository(db, config)}
}

package app

import (
	"context"
	"database/sql"
	"github.com/semelyanov86/notifications-runner/internal/config"
	"github.com/semelyanov86/notifications-runner/internal/handler"
	"github.com/semelyanov86/notifications-runner/internal/repository"
	"github.com/semelyanov86/notifications-runner/internal/service"
	"github.com/semelyanov86/notifications-runner/pkg/cache"
	"github.com/semelyanov86/notifications-runner/pkg/logger"
	"time"
)

// Run initializes whole application.
func Run(configPath string) {

	cfg := config.Init(configPath)
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(logger.ConvertErrorToStruct(err, 0, nil))
		return
	}
	memcache := cache.NewMemoryCache()

	repos := repository.NewRepositories(db, *cfg)
	services := service.NewServices(*cfg)
	handlers := handler.NewHandler(*repos, *services, memcache)
	err = handlers.Init()
	if err != nil {
		logger.Error(logger.ConvertErrorToStruct(err, 1025, nil))
		return
	}
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

package app

import (
	"Notifications/internal/config"
	"Notifications/internal/handler"
	"Notifications/internal/repository"
	"Notifications/internal/service"
	"Notifications/pkg/cache"
	"Notifications/pkg/logger"
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	handlers := handler.NewHandler(*repos, *services, memcache, db)
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

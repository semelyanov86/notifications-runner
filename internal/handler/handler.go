package handler

import (
	"Notifications/internal/domain"
	"Notifications/internal/repository"
	"Notifications/internal/service"
	"Notifications/pkg/cache"
	"Notifications/pkg/e"
	"Notifications/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var ErrNotSupportedType = errors.New("this type is not supported")

type Handler struct {
	Repos    repository.Repositories
	Services service.Services
	Cache    cache.Cache
	Db       *sql.DB
}

func NewHandler(repos repository.Repositories, service service.Services, cache cache.Cache, db *sql.DB) Handler {
	return Handler{
		Repos:    repos,
		Services: service,
		Cache:    cache,
		Db:       db,
	}
}

func (h Handler) Init() error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	lastJobID, err := h.Repos.Entities.GetLastProcessedEntity()
	if errors.Is(repository.ErrRecordNotFound, err) {
		lastJobID = 1
	} else if err != nil {
		cancel()
		return e.Wrap("can not get last not processed entry", err)
	}

	go func() {
		err = h.listenForNotifications(ctx, &wg, lastJobID)
		if err != nil {
			logger.Error(logger.LogMessage{
				Msg:        err.Error(),
				Code:       "1045",
				Properties: nil,
			})
		}
	}()
	logger.Info(logger.GenerateErrorMessageFromString("Started listening for new notifications, last notification was " + strconv.Itoa(lastJobID)))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	cancel()
	wg.Wait()

	if err := h.Db.Close(); err != nil {
		return err
	}
	return nil
}

func (h Handler) listenForNotifications(ctx context.Context, wg *sync.WaitGroup, lastJobID int) error {
	id := lastJobID
	for {
		select {
		case <-ctx.Done():
			// The context has been canceled, stop the goroutine
			return ctx.Err()
		default:
			time.Sleep(5 * time.Second) // Wait for 5 seconds before polling again

			entity, err := h.Repos.Entities.GetNextNotProcessedEntity(id)

			if errors.Is(repository.ErrRecordNotFound, err) {
				continue
			}
			if err != nil {
				logger.Error(logger.LogMessage{
					Msg:        err.Error(),
					Code:       "1045",
					Properties: nil,
				})
				continue
			}
			entity, err = h.Repos.Entities.GetEntityById(entity.Crmid)
			if err != nil {
				return e.Wrap("can not get record data", err)
			}
			id = entity.Crmid
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = h.doSend(entity)
				if err != nil {
					logger.Error(logger.LogMessage{
						Msg:        err.Error(),
						Code:       "1065",
						Properties: nil,
					})
				}
			}()
		}
	}
}

func (h Handler) doSend(entity *domain.Entity) error {
	switch entity.NotifyType {
	case "Slack":
		_, channelId, err := h.Services.Slack.Send(*entity)
		entity.ChatId.String = channelId
		entity.ChatId.Valid = true
		if err != nil {
			return h.Repos.Entities.MarkAsError(entity, err)
		}

		entity.NotifyDateTime = time.Now()
	default:
		return ErrNotSupportedType
	}
	return h.Repos.Entities.MarkAsSent(entity)
}

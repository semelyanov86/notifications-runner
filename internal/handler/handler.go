package handler

import (
	"github.com/semelyanov86/notifications-runner/internal/repository"
	"github.com/semelyanov86/notifications-runner/internal/service"
	"github.com/semelyanov86/notifications-runner/pkg/cache"
)

type Handler struct {
	Repos    repository.Repositories
	Services service.Services
	Cache    cache.Cache
}

func NewHandler(repos repository.Repositories, service service.Services, cache cache.Cache) Handler {
	return Handler{
		Repos:    repos,
		Services: service,
		Cache:    cache,
	}
}

func (h Handler) Init() error {
	return nil
}

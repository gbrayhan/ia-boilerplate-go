package handlers

import (
	"go.uber.org/zap"
	"ia-boilerplate/src/infrastructure"
	"ia-boilerplate/src/repository"
)

type Handler struct {
	Repository     *repository.Repository
	Infrastructure *infrastructure.Infrastructure
	Logger         *zap.Logger
}

func NewHandler(repository *repository.Repository, logger *zap.Logger, infrastructure *infrastructure.Infrastructure) *Handler {
	return &Handler{
		Repository:     repository,
		Logger:         logger,
		Infrastructure: infrastructure,
	}
}

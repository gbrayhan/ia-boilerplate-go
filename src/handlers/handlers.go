package handlers

import (
	"ia-boilerplate/src/infrastructure"
	"ia-boilerplate/src/repository"
)

type Handler struct {
	Repository *repository.Repository
	Auth       *infrastructure.Auth
	Logger     *infrastructure.Logger
}

func NewHandler(repository *repository.Repository, logger *infrastructure.Logger, auth *infrastructure.Auth) *Handler {
	return &Handler{
		Repository: repository,
		Logger:     logger,
		Auth:       auth,
	}
}

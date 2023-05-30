package handlers

import (
	"github.com/TheVitik/sport-together-go/internal/repositories"
)

type Handler struct {
	repository *repositories.Repository
}

func NewHandler(repository *repositories.Repository) *Handler {
	return &Handler{repository: repository}
}

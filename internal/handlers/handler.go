package handlers

import (
	"fmt"

	"github.com/epic55/BankApp/internal/models"
	"github.com/epic55/BankApp/internal/repository"
)

type Handler struct {
	R    *repository.Repository
	Cnfg *models.Config
}

func NewHandler(repo *repository.Repository, config *models.Config) *Handler {
	if repo == nil {
		fmt.Println("Failed to initialize the repo")
	}

	return &Handler{
		R:    repo,
		Cnfg: config,
	}
}

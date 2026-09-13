package handler

import (
	"net/http"

	"github.com/abdelbassat/github-stats-api/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {

}

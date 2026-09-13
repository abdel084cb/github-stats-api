package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/abdelbassat/github-stats-api/internal/apperrors"
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

func writeError(w http.ResponseWriter, err error) {
	httpStatus := http.StatusInternalServerError

	errorMap := map[string]string{
		"error": "internal server error",
	}

	switch {
	case errors.Is(err, apperrors.ErrInvalidUsername):
		httpStatus = http.StatusBadRequest
		errorMap["error"] = "invalid username"

	case errors.Is(err, apperrors.ErrUserNotFound):
		httpStatus = http.StatusNotFound
		errorMap["error"] = "user not found"

	case errors.Is(err, apperrors.ErrUpstreamUnavailable):
		httpStatus = http.StatusServiceUnavailable
		errorMap["error"] = "upstream unavailable"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	json.NewEncoder(w).Encode(errorMap) // unhandled error
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	stats, err := h.service.GetStatsByUsername(r.Context(), username)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

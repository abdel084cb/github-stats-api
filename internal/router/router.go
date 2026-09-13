package router

import (
	"net/http"

	"github.com/abdelbassat/github-stats-api/internal/handler"
)

func NewRouter(handler *handler.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats/{username}", handler.GetStats)
	return mux
}

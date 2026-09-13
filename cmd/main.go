package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/abdelbassat/github-stats-api/internal/githubclient"
	"github.com/abdelbassat/github-stats-api/internal/handler"
	"github.com/abdelbassat/github-stats-api/internal/repository"
	"github.com/abdelbassat/github-stats-api/internal/router"
	"github.com/abdelbassat/github-stats-api/internal/service"
	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite", "stats.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repositoryStats := repository.NewRepository(db)

	if err := repositoryStats.CreateDatabase(db); err != nil {
		log.Fatal(err)
	}

	gitHubHTTPClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	gitHubClient := githubclient.NewClient(gitHubHTTPClient)

	serviceStats := service.NewService(
		repositoryStats,
		gitHubClient,
	)

	handlerStats := handler.NewHandler(serviceStats)
	routerStats := router.NewRouter(handlerStats)

	server := &http.Server{
		Addr:              "localhost:8003",
		Handler:           routerStats,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server is running and listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

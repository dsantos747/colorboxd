package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	colorboxd "github.com/dsantos747/letterboxd_hue_sort/backend"
	"github.com/dsantos747/letterboxd_hue_sort/backend/redis"
)

func main() {
	if err := colorboxd.LoadEnv(); err != nil {
		slog.Warn("could not load environment variables from .env file", "err", err)
	}

	rc, err := redis.New(os.Getenv("REDIS_URL"))
	if err != nil {
		slog.Error("could not connect to redis", "err", err)
		os.Exit(1)
	}
	colorboxd.SetRedisClient(rc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/auth", colorboxd.AuthUser)
	mux.HandleFunc("GET /api/v1/lists", colorboxd.GetLists)
	mux.HandleFunc("GET /api/v1/sort", colorboxd.SortListById)
	mux.HandleFunc("POST /api/v1/write", colorboxd.WriteList)
	mux.HandleFunc("OPTIONS /api/v1/write", colorboxd.WriteList)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 110 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("starting server", "port", port)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}

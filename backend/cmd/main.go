package main

import (
	"log"
	"net/http"
	"os"
	"time"

	colorboxd "github.com/dsantos747/letterboxd_hue_sort/backend"
	"github.com/dsantos747/letterboxd_hue_sort/backend/redis"
)

func main() {
	if err := colorboxd.LoadEnv(); err != nil {
		log.Printf("Could not load environment variables from .env file: %v", err)
	}

	rc, err := redis.New(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
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

	log.Printf("Starting server on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

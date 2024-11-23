package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

var CR *ColorRepo

func main() {
	ctx := context.Background()

	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := LoadEnv()
	if err != nil {
		l.Error("failed to load envs", "err", err)
		return
	}

	CR, err = NewColorRepo(l)
	if err != nil {
		l.Error("failed to load colorRepo", "err", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("api/v2/AuthUser", AuthUser)
	mux.HandleFunc("api/v2/GetLists", GetLists)
	mux.HandleFunc("api/v2/SortList", SortList)
	mux.HandleFunc("api/v2/WriteList", WriteList)

	server := http.Server{
		Addr:    ":8080",
		Handler: headerMiddleware(mux),
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Log(ctx, slog.LevelError, "failed to listenandserve on main server", "err", err)
		os.Exit(1)
	}
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set necessary headers for CORS and cache policy
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Cache-Control", "private, max-age=3570") // Expire time of token (-30s for safety)

		next.ServeHTTP(w, r)
	})
}

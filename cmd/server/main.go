package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/simonecolaci/subito-assignment/internal/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	handler.NewHealthHandler().RegisterRoutes(mux)

	slog.Info("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}

package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/simonecolaci/subito-assignment/internal/handler"

	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		slog.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database", "err", err)
		}
	}()

	mux := http.NewServeMux()

	handler.NewHealthHandler().RegisterRoutes(mux)

	slog.Info("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}

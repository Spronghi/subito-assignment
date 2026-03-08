package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/simonecolaci/subito-assignment/internal/handler"
	"github.com/simonecolaci/subito-assignment/internal/service"

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

	productService := service.NewProductService()

	// TODO: handle graceful shutdown by listening to OS signals in a separate goroutine
	mux := http.NewServeMux()

	handler.NewHealthHandler().RegisterRoutes(mux)
	handler.NewProductHandler(productService).RegisterRoutes(mux)

	slog.Info("Starting server on :8080")

	// TODO: check for proper server configuration (e.g., timeouts)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Server failed to start", "err", err)
		os.Exit(1)
	}
}

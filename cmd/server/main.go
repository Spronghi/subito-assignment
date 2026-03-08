package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/simonecolaci/subito-assignment/internal/handler"
	"github.com/simonecolaci/subito-assignment/internal/repository"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	//NOTE: using in-memory SQLite for simplicity, in a real application we would use a persistent database
	db, err := repository.NewSQLiteDB(":memory:")

	if err != nil {
		slog.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database", "err", err)
		}
	}()

	productRepo, err := repository.NewSQLiteProductRepository(db)
	if err != nil {
		slog.Error("Failed to initialize product repository", "err", err)
		os.Exit(1)
	}

	if err := productRepo.Populate(); err != nil {
		slog.Error("Failed to populate product repository", "err", err)
		os.Exit(1)
	}

	productService := service.NewProductService(productRepo)

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

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/jha-captech/blog/internal/config"
	"github.com/jha-captech/blog/internal/database"
	"github.com/jha-captech/blog/internal/middleare"
	"github.com/jha-captech/blog/internal/routes"
	"github.com/jha-captech/blog/internal/services"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server encountered an error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Load and validate environment config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in main.run] failed to load config: %w", err)
	}

	// Create a structured logger, which will print logs in json format to the
	// writer we specify.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	// Create a new DB connection using environment config
	db, err := database.Connect(ctx, logger, fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUserName,
		cfg.DBUserPassword,
		cfg.DBName,
		cfg.DBPort,
	))
	if err != nil {
		return fmt.Errorf("[in main.run] failed to open database: %w", err)
	}
	defer func() {
		logger.DebugContext(ctx, "Closing database connection")
		if err = db.Close(); err != nil {
			logger.ErrorContext(ctx, "Failed to close database connection", "err", err)
		}
	}()
	logger.InfoContext(ctx, "Connected successfully to the database")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.CacheHost, cfg.CachePort),
		Password: cfg.CachePassword,
		DB:       cfg.CacheDB,
	})

	usersService := services.NewUsersService(logger, db, rdb)

	// Create a serve mux to act as our route multiplexer
	mux := http.NewServeMux()

	// Add our routes to the mux
	routes.AddRoutes(mux, logger, usersService, fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port))

	// Wrap the mux with middleware
	wrappedMux := middleare.Logger(logger)(mux)

	// Create a new http server with our mux as the handler
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: wrappedMux,
	}

	errChan := make(chan error)

	// Server run context
	ctx, done := context.WithCancel(ctx)
	defer done()

	// Handle graceful shutdown with go routine on SIGINT
	go func() {
		// create a channel to listen for SIGINT and then block until it is received
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig

		logger.DebugContext(ctx, "Received SIGINT, shutting down server")

		// Create a context with a timeout to allow the server to shut down gracefully
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Shutdown the server. If an error occurs, send it to the error channel
		if err = httpServer.Shutdown(ctx); err != nil {
			errChan <- fmt.Errorf("[in main.run] failed to shutdown http server: %w", err)
			return
		}

		// Close the idle connections channel, unblocking `run()`
		done()
	}()

	// Start the http server
	logger.InfoContext(ctx, "listening", slog.String("address", httpServer.Addr))
	if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		switch {
		// once httpServer.Shutdown is called, it will always return an
		// http.ErrServerClosed error and we don't care about that error so we will
		// break.
		case errors.Is(err, http.ErrServerClosed):
			break
		default:
			return fmt.Errorf("[in main.run] failed to listen and serve: %w", err)
		}
	}

	// block until the server is shut down or an error occurs
	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}

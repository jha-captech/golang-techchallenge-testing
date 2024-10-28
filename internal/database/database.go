package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect establishes a database connection and returns the connection.
func Connect(ctx context.Context, logger *slog.Logger, connectionString string) (*sql.DB, error) {
	// Create a new DB connection using environment config
	logger.DebugContext(ctx, "Connecting to database")
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("[in database.Connect] failed to open database: %w", err)
	}

	// Ping the database to verify connection
	logger.DebugContext(ctx, "Pinging database")
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("[in database.Connect] failed to ping database: %w", err)
	}

	return db, nil
}

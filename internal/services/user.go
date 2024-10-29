package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jha-captech/blog/internal/models"
	"github.com/redis/go-redis/v9"
)

// UsersService is a service capable of performing CRUD operations for
// models.User models.
type UsersService struct {
	logger *slog.Logger
	db     *sql.DB
	cache  *Client
}

// NewUsersService creates a new UsersService and returns a pointer to it.
func NewUsersService(logger *slog.Logger, db *sql.DB, rdb *redis.Client, expiration time.Duration) *UsersService {
	return &UsersService{
		logger: logger,
		db:     db,
		cache: &Client{
			Client:     rdb,
			expiration: expiration,
		},
	}
}

// CreateUser attempts to create the provided user, returning a fully hydrated
// models.User or an error.
func (s *UsersService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return models.User{}, nil
}

// ReadUser attempts to read a user from the database using the provided id. A
// fully hydrated models.User or error is returned.
func (s *UsersService) ReadUser(ctx context.Context, id uint64) (models.User, error) {
	logger := s.logger.With(slog.String("func", "services.UsersService.ReadUser"))
	logger.DebugContext(ctx, "Getting user", "id", id)

	// Check the cache for the user object
	logger.DebugContext(ctx, "Reading user from cache", "id", id)

	var user models.User
	found, err := s.cache.Get(ctx, strconv.FormatUint(id, 10)).UnmarshalJSON(&user)
	if err != nil {
		return models.User{}, fmt.Errorf(
			"[in services.UsersService.ReadUser] failed to read user from cache: %w",
			err,
		)
	}

	// If the user was found in the cache, return it
	if found {
		return user, nil
	}

	// If the user was not found in the cache, read it from the database
	logger.DebugContext(ctx, "Reading user from database", "id", id)
	row := s.db.QueryRowContext(
		ctx,
		`
		SELECT id,
		       name,
		       email,
		       password
		FROM users
		WHERE id = $1::int
		`,
		id,
	)

	// Scan the row into the user object
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, nil
		default:
			return models.User{}, fmt.Errorf(
				"[in services.UsersService.ReadUser] failed to read user: %w",
				err,
			)
		}
	}

	// Write the user to the cache
	logger.DebugContext(ctx, "Setting user in cache", "id", id)
	if err = s.cache.SetJSON(ctx, strconv.FormatUint(id, 10), user); err != nil {
		return models.User{}, fmt.Errorf(
			"[in services.UsersService.ReadUser] failed to write user to cache: %w",
			err,
		)
	}

	return user, nil
}

// UpdateUser attempts to perform an update of the user with the provided id,
// updating, it to reflect the properties on the provided patch object. A
// models.User or an error.
func (s *UsersService) UpdateUser(ctx context.Context, id uint64, patch models.User) (models.User, error) {
	return models.User{}, nil
}

// DeleteUser attempts to delete the user with the provided id. An error is
// returned if the delete fails.
func (s *UsersService) DeleteUser(ctx context.Context, id uint64) error {
	return nil
}

// ListUsers attempts to list all users in the database. A slice of models.User
// or an error is returned.
func (s *UsersService) ListUsers(ctx context.Context, id uint64) ([]models.User, error) {
	return []models.User{}, nil
}

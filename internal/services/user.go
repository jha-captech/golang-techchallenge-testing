package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jha-captech/blog/internal/models"
)

// UsersService is a service capable of performing CRUD operations for
// models.User models.
type UsersService struct {
	logger *slog.Logger
	db     *sql.DB
}

// NewUsersService creates a new UsersService and returns a pointer to it.
func NewUsersService(logger *slog.Logger, db *sql.DB) *UsersService {
	return &UsersService{
		logger: logger,
		db:     db,
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
	s.logger.DebugContext(ctx, "Reading user", "id", id)

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

	var user models.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
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

	return user, nil
}

// UpdateUser attempts to perform an update of the user with the provided id,
// updating, it to reflect the properties on the provided patch object. A
// models.User or an error.
func (s *UsersService) UpdateUser(ctx context.Context, id uint64, patch models.User) (models.User, error) {
	return models.User{}, nil
}

// CreateUser attempts to create the provided user, returning a fully hydrated
// models.User or an error.
func (s *UsersService) DeleteUser(ctx context.Context, id uint64) error {
	return nil
}

// CreateUser attempts to create the provided user, returning a fully hydrated
// models.User or an error.
func (s *UsersService) ListUsers(ctx context.Context, id uint64) ([]models.User, error) {
	return []models.User{}, nil
}

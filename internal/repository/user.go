package repository

import (
	"context"
	"database/sql"

	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// UserRepository handles user database operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	).Scan(&user.ID)
	if err != nil {
		logger.Error("Failed to create user", err, "email", user.Email)
		return err
	}

	logger.Info("User created", "email", user.Email, "user_id", user.ID)
	return nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password FROM users WHERE email = $1 LIMIT 1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Error("Failed to find user", err, "email", email)
		return nil, err
	}

	return &user, nil
}

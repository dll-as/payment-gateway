package repository

import (
	"context"
	"database/sql"

	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// UserRepository handles user database operations
type UserRepository interface {
	Register(ctx context.Context, user *models.User) error
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
func (r *userRepository) Register(ctx context.Context, user *models.User) error {
	if err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (email, password, full_name, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Email, user.Password, user.FullName, user.Role,
	).Scan(&user.ID); err != nil {
		logger.Error("Failed to create user", err, "email", user.Email)
		return err
	}

	logger.Info("User created", "email", user.Email, "user_id", user.ID)
	return nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password FROM users WHERE email = $1 LIMIT 1", email,
	).Scan(&user.ID, &user.Email, &user.Password); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		logger.Error("Failed to find user", err, "email", email)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUser(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	if err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password, full_name, role FROM users WHERE email = $1 LIMIT 1", id,
	).Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.Role); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		logger.Error("Failed to find user", err, "id", id)
		return nil, err
	}

	return &user, nil
}

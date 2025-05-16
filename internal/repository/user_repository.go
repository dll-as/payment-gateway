package repository

// package repositories

// import (
// 	"context"
// 	"database/sql"
// 	"errors"

// 	"github.com/rezatg/payment-gateway/internal/models"
// )

// type UserRepository interface {
// 	Create(ctx context.Context, user *models.User) error
// 	FindByEmail(ctx context.Context, email string) (*models.User, error)
// }

// type userRepository struct {
// 	db *sql.DB
// }

// func NewUserRepository(db *sql.DB) UserRepository {
// 	return &userRepository{db}
// }

// func (r *userRepository) Create(ctx context.Context, user *models.User) error {
// 	if err := r.db.QueryRowContext(ctx,
// 		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
// 		user.Email, user.Password,
// 	).Scan(&user.ID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
// 	var user models.User
// 	if err := r.db.QueryRowContext(
// 		ctx, "SELECT id, email, password FROM users WHERE email = $1 LIMIT 1", email,
// 	).Scan(&user.ID, &user.Email, &user.Password); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, nil
// 		}

// 		return nil, err
// 	}

// 	return &user, nil
// }

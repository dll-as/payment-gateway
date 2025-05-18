// package auth

// import (
// 	"context"
// 	"errors"
// 	"time"

// 	"github.com/rezatg/payment-gateway/internal/models"
// 	"github.com/rezatg/payment-gateway/internal/repositories"
// 	errs "github.com/rezatg/payment-gateway/pkg/errors"
// 	"github.com/rezatg/payment-gateway/pkg/logger"
// 	"github.com/rezatg/payment-gateway/pkg/utils"
// )

// // AuthService defines the interface for authentication operations
// type AuthService interface {
// 	Register(ctx context.Context, email, password string) error
// 	Login(ctx context.Context, email, password string) (string, error)
// 	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
// 	ValidateToken(ctx context.Context, token string) (string, error)
// }

// type authService struct {
// 	userRepo repositories.UserRepository
// }

// func New(userRepo repositories.UserRepository) AuthService {
// 	return &authService{userRepo}
// }

// func (s *authService) Register(ctx context.Context, email, password string) error {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	user, err := s.userRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		return err
// 	} else if user != nil {
// 		return errors.New("email already exists")
// 	}

// 	hashedPassword, err := utils.HashPassword(password)
// 	if err != nil {
// 		return err
// 	}

// 	newUser := &models.User{
// 		Email:    email,
// 		Password: hashedPassword,
// 	}
// 	return s.userRepo.Create(ctx, newUser)
// }

// func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
// 	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	user, err := s.FindUserByEmail(timeoutCtx, email)
// 	if err != nil {
// 		return "", err
// 	} else if user == nil {
// 		return "", errors.New("user not found")
// 	}

// 	if !utils.CheckPasswordHash(password, user.Password) {
// 		return "", errors.New("invalid credentials")
// 	}

// 	token, err := generateJWTToken(user.ID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// func (s *authService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
// 	return s.userRepo.FindByEmail(ctx, email)
// }

// // ValidateToken validates a JWT token and returns the user ID
// func (s *authService) ValidateToken(ctx context.Context, token string) (string, error) {
// 	_, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	userID, err := validateJWTToken(token)
// 	if err != nil {
// 		logger.Error("Failed to validate JWT token", err, "token", token[:10]+"...")
// 		return "", errs.NewUnauthorizedError("Invalid token", err.Error())
// 	}

// 	logger.Debug("Token validated successfully", "user_id", userID)
// 	return userID, nil
// }

package auth

import (
	"context"
	"time"

	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/internal/repository"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
	"github.com/rezatg/payment-gateway/pkg/utils"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

// authService implements AuthService
type authService struct {
	userRepo repository.UserRepository
}

// New creates a new AuthService
func New(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, email, password string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Error("Failed to check existing user", err, "email", email)
		return errors.NewInternalServerError("Failed to register", err.Error())
	} else if user != nil {
		return errors.NewConflictError("Email already exists")
	}

	password, err = utils.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password", err, "email", email)
		return errors.NewInternalServerError("Failed to register", err.Error())
	}

	newUser := &models.User{
		Email:    email,
		Password: password,
		FullName: "reza",
		Role:     0,
	}
	if err = s.userRepo.Register(ctx, newUser); err != nil {
		logger.Error("Failed to create user", err, "email", email)
		return errors.NewInternalServerError("Failed to register", err.Error())
	}

	logger.Info("User registered", "email", email, "user_id", newUser.ID)
	return nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Error("Failed to find user", err, "email", email)
		return "", errors.NewInternalServerError("Failed to login", err.Error())
	}
	if user == nil {
		return "", errors.NewNotFoundError("User not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.NewUnauthorizedError("Invalid credentials")
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate JWT token", err, "user_id", user.ID)
		return "", errors.NewInternalServerError("Failed to generate token", err.Error())
	}

	logger.Info("User logged in", "email", email, "user_id", user.ID)
	return token, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *authService) ValidateToken(_ context.Context, token string) (string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	userID, err := validateJWTToken(token)
	if err != nil {
		logger.Error("Failed to validate JWT token", err, "token", token[:10]+"...")
		return "", errors.NewUnauthorizedError("Invalid token")
	}

	logger.Debug("Token validated", "user_id", userID)
	return userID, nil
}

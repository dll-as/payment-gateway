package auth

import (
	"context"
	"errors"
	"time"

	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(ctx context.Context, email, password string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	} else if user != nil {
		return errors.New("email already exists")
	}

	newUser := &models.User{
		Email:    email,
		Password: password,
	}
	return s.userRepo.Create(ctx, newUser)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.FindUserByEmail(timeoutCtx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	if user.Password != password {
		return "", errors.New("invalid credentials")
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

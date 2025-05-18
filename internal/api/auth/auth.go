package auth

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/internal/services/auth"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	service   auth.AuthService
	validator *validator.Validate
}

// New creates a new AuthHandler
func New(service auth.AuthService) *AuthHandler {
	return &AuthHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Login handles user login
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req models.LoginRequest
	contentType := c.Get("Content-Type")
	var err error
	if strings.Contains(contentType, "application/json") {
		fmt.Println("dd")
		err = c.Bind().JSON(&req)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err = c.Bind().Form(&req)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Unsupported Content-Type",
			Code:    fiber.StatusBadRequest,
			Details: "Use application/json or application/x-www-form-urlencoded",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation failed",
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
		})
	}

	token, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		logger.Error("Login failed", err, "email", req.Email)
		return errors.HandleAPIError(c, err)
	}

	logger.Info("User logged in", "email", req.Email)
	return c.JSON(models.LoginResponse{Token: token})
}

// Register handles user registration
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation failed",
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
		})
	}

	if err := h.service.Register(c.Context(), req.Email, req.Password); err != nil {
		logger.Error("Registration failed", err, "email", req.Email)
		return errors.HandleAPIError(c, err)
	}

	logger.Info("User registered", "email", req.Email)
	return c.Status(fiber.StatusCreated).JSON(map[string]string{"message": "User registered successfully"})
}

// GetAuthService provides auth validation for middleware
func (h *AuthHandler) GetAuthService() auth.AuthService {
	return h.service
}

package controllers

import (
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/app/util/jwt"
	"ayo-baca-buku/app/util/logger"
	"ayo-baca-buku/app/util/validation"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthController struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewAuthController(DB *gorm.DB) *AuthController {
	return &AuthController{
		DB:       DB,
		Validate: validator.New(),
	}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,alphanum"`
	Password string `json:"password" validate:"required,min=6,alphanum"`
}

type LoginResponse struct {
	Message string            `json:"message"`
	Token   string            `json:"token,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags Login
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} LoginResponse
// @Router /login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	logger.Info("AuthController.Login Begin")

	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(LoginResponse{
			Message: "Invalid Request",
			Errors: map[string]string{
				"body": "Failed to parse request body",
			},
		})
	}

	if err := c.Validate.Struct(&req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		validationErrors := make(map[string]string)

		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(LoginResponse{
			Message: "Validation failed",
			Errors:  validationErrors,
		})
	}

	user := models.User{}
	if err := c.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		logger.Error("User not found", zap.Error(err))
		return ctx.Status(fiber.StatusNotFound).JSON(LoginResponse{
			Message: "Invalid credentials (1)",
		})
	}

	if user.DeletedBy != 0 {
		logger.Error("User deleted")
		return ctx.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
			Message: "User deleted (soft)",
		})
	}

	if !jwt.CheckPasswordHash(req.Password, user.Password) {
		logger.Error("Invalid password")
		return ctx.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
			Message: "Invalid credentials (2)",
		})
	}

	token, err := jwt.GenerateToken(user.UID, user.Username)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(LoginResponse{
			Message: "Failed to generate token",
		})
	}

	if err := c.DB.Model(&user).Where("username = ?", user.Username).Update("token", token).Error; err != nil {
		logger.Error("Failed to update token", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(LoginResponse{
			Message: "Failed to update token",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(LoginResponse{
		Message: "Success",
		Token:   token,
	})
}

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required"`
	Username             string `json:"username" validate:"required,alphanum,unique_username"`
	Email                string `json:"email" validate:"required,email,unique_email"`
	Password             string `json:"password" validate:"required,alphanum,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password,min=6"`
}

type RegisterResponse struct {
	Message string            `json:"message"`
	Data    *models.User      `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Register godoc
// @Summary Register
// @Description Register
// @Tags Register
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} RegisterResponse
// @Router /register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	logger.Info("AuthController.Register Begin")

	c.Validate.RegisterValidation("unique_username", validation.UniqueUsername(c.DB, 0))
	c.Validate.RegisterValidation("unique_email", validation.UniqueEmail(c.DB, 0))

	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
			Message: "Invalid Request",
			Errors: map[string]string{
				"body": "Failed to parse request body",
			},
		})
	}

	if err := c.Validate.Struct(&req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		validationErrors := make(map[string]string)

		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
			Message: "Validation failed",
			Errors:  validationErrors,
		})
	}

	hash, err := jwt.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
			Message: "Failed to hash password",
		})
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: hash,
		Role:     "user",
	}

	if err := c.DB.Create(&user).Error; err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
			Message: "Failed to create user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(RegisterResponse{
		Message: "Success",
		Data:    &user,
	})
}

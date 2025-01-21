package controllers

import (
	"ayo-baca-buku/app/util/logger"
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

	return ctx.Status(fiber.StatusOK).JSON(LoginResponse{
		Message: "Success",
	})
}

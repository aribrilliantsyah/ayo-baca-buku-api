package controllers

import (
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/app/util/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) *UserController {
	return &UserController{
		DB: DB,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	logger := logger.GetLogger() // Global logger instance

	logger.Info("Fetching all users")
	var users []*models.User
	if err := c.DB.Debug().Find(&users).Error; err != nil {
		logger.Error("Failed to fetch users", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch users",
		})
	}

	logger.Info("Fetched users successfully")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    users,
	})
}

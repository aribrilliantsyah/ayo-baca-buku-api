package controllers

import (
	"ayo-baca-buku/app/models"

	"github.com/gofiber/fiber/v2"
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
	var users []*models.User
	c.DB.Debug().Find(&users)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    users,
	})
}

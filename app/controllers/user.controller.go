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
	if err := c.DB.Find(&users).Error; err != nil {
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

// GetUserById godoc
// @Summary Get user by ID
// @Description Get a single user by their ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} fiber.Map{message=string, data=models.User}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /users/{id} [get]
func (c *UserController) GetUserById(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	userID := ctx.Params("id")

	logger.Info("Fetching user by ID", zap.String("userID", userID))

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found", zap.String("userID", userID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		logger.Error("Failed to fetch user", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user",
		})
	}

	logger.Info("User fetched successfully", zap.String("userID", userID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User Create Payload"
// @Success 201 {object} fiber.Map{message=string, data=models.User}
// @Failure 400 {object} fiber.Map{message=string, errors=map[string]string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	logger.Info("UserController.CreateUser Begin")

	// Initialize validator and register custom validations
	validate := validator.New()
	validate.RegisterValidation("unique_username", validation.UniqueUsername(c.DB, 0))
	validate.RegisterValidation("unique_email", validation.UniqueEmail(c.DB, 0))

	var req models.UserCreateRequest // Assuming UserCreateRequest is defined in models package
	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
			"errors": map[string]string{
				"body": "Failed to parse request body",
			},
		})
	}

	if err := validate.Struct(&req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		validationErrors := make(map[string]string)
		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	hash, err := jwt.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: hash,
		Role:     "user", // Default role, or get from request if applicable
	}

	if err := c.DB.Create(&user).Error; err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	// It's good practice to not return the password hash in the response.
	// Create a UserResponse struct or selectively pick fields.
	// For simplicity, returning the user object (excluding password).
	user.Password = "" // Clear password before sending response

	logger.Info("User created successfully", zap.Uint("userID", user.UID))
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user with the input payload
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UserUpdateRequest true "User Update Payload"
// @Success 200 {object} fiber.Map{message=string, data=models.User}
// @Failure 400 {object} fiber.Map{message=string, errors=map[string]string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	userID := ctx.Params("id")
	logger.Info("UserController.UpdateUser Begin", zap.String("userID", userID))

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found for update", zap.String("userID", userID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		logger.Error("Failed to fetch user for update", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user"})
	}

	// Initialize validator and register custom validations
	validate := validator.New()
	// Pass user.ID to ignore current user's email/username in unique checks
	validate.RegisterValidation("unique_username", validation.UniqueUsername(c.DB, user.ID))
	validate.RegisterValidation("unique_email", validation.UniqueEmail(c.DB, user.ID))

	var req models.UserUpdateRequest // Assuming UserUpdateRequest is defined in models package
	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body for update", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
			"errors":  map[string]string{"body": "Failed to parse request body"},
		})
	}

	if err := validate.Struct(&req); err != nil {
		logger.Error("Validation failed for update", zap.Error(err))
		validationErrors := make(map[string]string)
		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	// Update password if provided and valid
	if req.Password != "" {
		if req.PasswordConfirmation == "" || req.Password != req.PasswordConfirmation {
			logger.Warn("Password confirmation does not match for update", zap.String("userID", userID))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  map[string]string{"password_confirmation": "Password confirmation does not match"},
			})
		}
		hashedPassword, err := jwt.HashPassword(req.Password)
		if err != nil {
			logger.Error("Failed to hash new password for update", zap.Error(err), zap.String("userID", userID))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update password"})
		}
		user.Password = hashedPassword
	}

	// Save updates
	if err := c.DB.Save(&user).Error; err != nil {
		logger.Error("Failed to update user", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update user"})
	}

	user.Password = "" // Clear password before sending response
	logger.Info("User updated successfully", zap.String("userID", userID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    user,
	})
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	userID := ctx.Params("id")
	logger.Info("UserController.DeleteUser Begin", zap.String("userID", userID))

	var user models.User
	// First, check if the user exists
	if err := c.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found for hard delete", zap.String("userID", userID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		logger.Error("Failed to fetch user for hard delete", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user"})
	}

	// Perform hard delete
	// Use Unscoped to permanently delete the record, bypassing GORM's soft delete
	if err := c.DB.Unscoped().Delete(&models.User{}, userID).Error; err != nil {
		logger.Error("Failed to hard delete user", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete user"})
	}

	logger.Info("User hard deleted successfully", zap.String("userID", userID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

// SoftDeleteUser godoc
// @Summary Soft delete a user
// @Description Soft delete a user by their ID (sets DeletedAt and DeletedBy fields)
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /users/{id}/soft-delete [patch]
func (c *UserController) SoftDeleteUser(ctx *fiber.Ctx) error {
	logger := logger.GetLogger()
	userID := ctx.Params("id") // ID of the user to be soft-deleted
	logger.Info("UserController.SoftDeleteUser Begin", zap.String("userID", userID))

	// In a real application, you would get the ID of the user performing the action
	// from JWT claims or session. For now, let's use a placeholder.
	// adminUserID := jwt.GetUserIDFromToken(ctx) // Example of how you might get it
	adminUserID := int64(1) // Placeholder admin/system user ID

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found for soft delete", zap.String("userID", userID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		logger.Error("Failed to fetch user for soft delete", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user"})
	}

	// Update DeletedBy and then perform GORM's soft delete
	// GORM's Delete method will automatically set DeletedAt if the model has gorm.DeletedAt field
	if err := c.DB.Model(&user).Update("DeletedBy", adminUserID).Error; err != nil {
		logger.Error("Failed to update DeletedBy for soft delete", zap.Error(err), zap.String("userID", userID))
		// Proceed with soft delete even if DeletedBy update fails, or handle as critical error
	}

	if err := c.DB.Delete(&user).Error; err != nil {
		logger.Error("Failed to soft delete user", zap.Error(err), zap.String("userID", userID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to soft delete user"})
	}

	logger.Info("User soft deleted successfully", zap.String("userID", userID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User soft deleted successfully"})
}

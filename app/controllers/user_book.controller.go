package controllers

import (
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/app/util/logger"
	"ayo-baca-buku/app/util/validation" // Assuming validation utilities might be needed
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserBookController struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserBookController(DB *gorm.DB) *UserBookController {
	return &UserBookController{
		DB:       DB,
		Validate: validator.New(),
	}
}

// CreateUserBook godoc
// @Summary Create a new user book entry
// @Description Add a new book to a user's reading list.
// @Tags UserBook
// @Accept json
// @Produce json
// @Param userbook body models.UserBookCreateRequest true "User Book Create Payload"
// @Success 201 {object} fiber.Map{message=string, data=models.UserBook}
// @Failure 400 {object} fiber.Map{message=string, errors=map[string]string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /userbooks [post]
func (c *UserBookController) CreateUserBook(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	log.Info("UserBookController.CreateUserBook Begin")

	var req models.UserBookCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
			"errors":  map[string]string{"body": "Failed to parse request body"},
		})
	}

	// TODO: In a real app, UserID might come from JWT token/auth context.
	// For now, it's in the request. We should validate if this user exists.
	var user models.User
	if err := c.DB.First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("User not found for UserBook creation", zap.Uint("userID", req.UserID))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  map[string]string{"user_id": "User not found"},
			})
		}
		log.Error("Failed to check user existence", zap.Error(err), zap.Uint("userID", req.UserID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error checking user"})
	}

	// Initialize validator (could be part of controller struct if reused often without per-handler registration)
	// c.Validate.RegisterValidation("custom_validation_if_any", validation.CustomValidationFunction(c.DB))
	if err := c.Validate.Struct(&req); err != nil {
		log.Error("Validation failed for UserBook creation", zap.Error(err))
		validationErrors := make(map[string]string)
		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	userBook := models.UserBook{
		UserID:         req.UserID,
		Title:          req.Title,
		Author:         req.Author,
		Publisher:      req.Publisher,
		Cover:          req.Cover,
		TotalPages:     req.TotalPages,
		MotivationRead: req.MotivationRead,
		Status:         "reading", // Default status
		StartDate:      req.StartDate,
		// EndDate will be null initially
		CurrentPage: 0, // Default current page
		CreatedBy:   int64(req.UserID), // Placeholder for actor ID
		UpdatedBy:   int64(req.UserID), // Placeholder for actor ID
	}

	if err := c.DB.Create(&userBook).Error; err != nil {
		log.Error("Failed to create UserBook in database", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user book entry",
		})
	}

	log.Info("UserBook created successfully", zap.Uint("userBookID", userBook.ID))
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User book entry created successfully",
		"data":    userBook,
	})
}

// UpdateUserBook godoc
// @Summary Update an existing user book
// @Description Update details of an existing user book by its ID.
// @Tags UserBook
// @Accept json
// @Produce json
// @Param id path int true "UserBook ID"
// @Param userbook_update body models.UserBookUpdateRequest true "User Book Update Payload"
// @Success 200 {object} fiber.Map{message=string, data=models.UserBook}
// @Failure 400 {object} fiber.Map{message=string, errors=map[string]string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /userbooks/{id} [put]
func (c *UserBookController) UpdateUserBook(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	userBookID := ctx.Params("id")
	log.Info("UserBookController.UpdateUserBook Begin", zap.String("userBookID", userBookID))

	var req models.UserBookUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body for UserBook update", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
			"errors":  map[string]string{"body": "Failed to parse request body"},
		})
	}

	// Validate the request payload
	if err := c.Validate.Struct(&req); err != nil {
		log.Error("Validation failed for UserBook update", zap.Error(err))
		validationErrors := make(map[string]string)
		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	var userBook models.UserBook
	if err := c.DB.First(&userBook, userBookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("UserBook not found for update", zap.String("userBookID", userBookID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User book not found"})
		}
		log.Error("Failed to fetch UserBook for update", zap.Error(err), zap.String("userBookID", userBookID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user book"})
	}

	// TODO: Authorization check: Does the authenticated user own this book or have rights to update it?
	// For example: if authenticatedUserID != userBook.UserID && !isAdmin(authenticatedUser) { return ctx.Status(fiber.StatusForbidden).JSON(...) }
	// For now, assume actorID for UpdatedBy can be derived or is a placeholder
	actorID := userBook.UserID // Placeholder for who is performing the update

	// Apply updates from request
	if req.Title != "" {
		userBook.Title = req.Title
	}
	if req.Author != "" {
		userBook.Author = req.Author
	}
	if req.Publisher != "" { // omitempty means empty string is a valid "not provided"
		userBook.Publisher = req.Publisher
	}
	if req.Cover != "" { // omitempty means empty string is a valid "not provided"
		userBook.Cover = req.Cover
	}
	if req.TotalPages != nil {
		userBook.TotalPages = *req.TotalPages
	}
	if req.CurrentPage != nil {
		userBook.CurrentPage = *req.CurrentPage
	}
	if req.MotivationRead != "" { // omitempty means empty string is a valid "not provided"
		userBook.MotivationRead = req.MotivationRead
	}
	if req.Status != "" {
		userBook.Status = req.Status
	}
	if !req.StartDate.IsZero() { // Check if StartDate is provided (not its zero value)
		userBook.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() { // Check if EndDate is provided
		userBook.EndDate = req.EndDate
	} else if req.Status == "finished" && userBook.EndDate.IsZero() {
		// If status is marked 'finished' and EndDate was not explicitly set, set it to now.
		userBook.EndDate = time.Now()
	}


	userBook.UpdatedBy = int64(actorID) // Placeholder

	if err := c.DB.Save(&userBook).Error; err != nil {
		log.Error("Failed to update UserBook in database", zap.Error(err), zap.Uint("userBookID", userBook.ID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user book entry",
		})
	}

	log.Info("UserBook updated successfully", zap.Uint("userBookID", userBook.ID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User book entry updated successfully",
		"data":    userBook,
	})
}

// DeleteUserBook godoc
// @Summary Soft delete a user book
// @Description Soft delete a user book by its ID.
// @Tags UserBook
// @Accept json
// @Produce json
// @Param id path int true "UserBook ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 403 {object} fiber.Map{message=string} // Forbidden if user doesn't own the book
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /userbooks/{id} [delete]
func (c *UserBookController) DeleteUserBook(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	userBookID := ctx.Params("id")
	log.Info("UserBookController.DeleteUserBook Begin", zap.String("userBookID", userBookID))

	var userBook models.UserBook
	if err := c.DB.First(&userBook, userBookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("UserBook not found for deletion", zap.String("userBookID", userBookID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User book not found"})
		}
		log.Error("Failed to fetch UserBook for deletion", zap.Error(err), zap.String("userBookID", userBookID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user book"})
	}

	// TODO: Authorization check: Does the authenticated user own this book or have rights to delete it?
	// For example:
	// authenticatedUserID := GetAuthenticatedUserID(ctx) // Implement this helper
	// if authenticatedUserID != userBook.UserID && !IsAdmin(authenticatedUserID) {
	//     log.Warn("User not authorized to delete UserBook", zap.Uint("userBookID", userBook.ID), zap.Uint("attemptedByUserID", authenticatedUserID))
	//     return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You are not authorized to delete this book entry"})
	// }
	actorID := userBook.UserID // Placeholder for who is performing the delete

	// Update DeletedBy before soft deleting
	// GORM's Delete will set DeletedAt automatically
	if err := c.DB.Model(&userBook).Update("DeletedBy", int64(actorID)).Error; err != nil {
		// Log the error but proceed with delete, as setting DeletedBy is audit info.
		// Depending on requirements, this could be a critical failure.
		log.Error("Failed to update DeletedBy for UserBook soft delete", zap.Error(err), zap.Uint("userBookID", userBook.ID))
	}

	if err := c.DB.Delete(&userBook).Error; err != nil {
		log.Error("Failed to soft delete UserBook in database", zap.Error(err), zap.Uint("userBookID", userBook.ID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user book entry",
		})
	}

	log.Info("UserBook soft deleted successfully", zap.Uint("userBookID", userBook.ID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User book entry deleted successfully"})
}

// GetAllUserBooks godoc
// @Summary Get all user books
// @Description Get a list of all user books, optionally filtered by user_id.
// @Tags UserBook
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by User ID"
// @Success 200 {object} fiber.Map{message=string, data=[]models.UserBook}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /userbooks [get]
func (c *UserBookController) GetAllUserBooks(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	log.Info("UserBookController.GetAllUserBooks Begin")

	var userBooks []models.UserBook
	query := c.DB

	// Optional filtering by user_id
	userID := ctx.QueryInt("user_id")
	if userID > 0 {
		log.Info("Filtering UserBooks by UserID", zap.Int("userID", userID))
		query = query.Where("user_id = ?", userID)
	}

	// Preload ReadingActivities and User for richer data. Consider if this is always needed.
	// For now, let's preload User to show who the book belongs to if not filtering.
	// ReadingActivities might be too much for a general list, better for GetUserBookByID.
	if err := query.Preload("User").Find(&userBooks).Error; err != nil {
		log.Error("Failed to fetch user books from database", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user books",
		})
	}

	log.Info("UserBooks fetched successfully", zap.Int("count", len(userBooks)))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User books fetched successfully",
		"data":    userBooks,
	})
}

// GetUserBookByID godoc
// @Summary Get a user book by its ID
// @Description Get details of a specific user book, including its owner and reading activities.
// @Tags UserBook
// @Accept json
// @Produce json
// @Param id path int true "UserBook ID"
// @Success 200 {object} fiber.Map{message=string, data=models.UserBook}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /userbooks/{id} [get]
func (c *UserBookController) GetUserBookByID(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	userBookID := ctx.Params("id") // This will be a string, needs conversion if your ID is int
	log.Info("UserBookController.GetUserBookByID Begin", zap.String("userBookID", userBookID))

	var userBook models.UserBook

	// Preload User and ReadingActivities for detailed view
	// Convert userBookID to appropriate type for GORM if necessary (e.g., to uint)
	// For now, GORM might handle string-to-int conversion for primary keys, but being explicit is better.
	// Let's assume ID in path is parseable to uint for the model's ID type.
	if err := c.DB.Preload("User").Preload("ReadingActivities").First(&userBook, userBookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("UserBook not found by ID", zap.String("userBookID", userBookID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User book not found",
			})
		}
		log.Error("Failed to fetch UserBook by ID from database", zap.Error(err), zap.String("userBookID", userBookID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user book",
		})
	}

	// TODO: Authorization check: Does the authenticated user own this book or have rights to view it?
	// For example: if authenticatedUserID != userBook.UserID && !isAdmin(authenticatedUser) { return ctx.Status(fiber.StatusForbidden).JSON(...) }

	log.Info("UserBook fetched successfully by ID", zap.Uint("userBookID", userBook.ID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User book fetched successfully",
		"data":    userBook,
	})
}

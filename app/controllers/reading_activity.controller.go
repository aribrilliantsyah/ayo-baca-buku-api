package controllers

import (
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/app/util/logger"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReadingActivityController struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewReadingActivityController(DB *gorm.DB) *ReadingActivityController {
	return &ReadingActivityController{
		DB:       DB,
		Validate: validator.New(),
	}
}

// CreateReadingActivity godoc
// @Summary Create a new reading activity
// @Description Add a new reading activity for a user's book.
// @Tags ReadingActivity
// @Accept json
// @Produce json
// @Param reading_activity body models.ReadingActivityCreateRequest true "Reading Activity Create Payload"
// @Success 201 {object} fiber.Map{message=string, data=models.ReadingActivity}
// @Failure 400 {object} fiber.Map{message=string, errors=map[string]string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string}
// @Router /reading-activities [post]
func (c *ReadingActivityController) CreateReadingActivity(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	log.Info("ReadingActivityController.CreateReadingActivity Begin")

	var req models.ReadingActivityCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
			"errors":  map[string]string{"body": "Failed to parse request body"},
		})
	}

	if err := c.Validate.Struct(&req); err != nil {
		log.Error("Validation failed for ReadingActivity creation", zap.Error(err))
		validationErrors := make(map[string]string)
		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	// Verify the UserBook exists
	var userBook models.UserBook
	if err := c.DB.First(&userBook, req.UserBookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("UserBook not found for ReadingActivity creation", zap.Uint("userBookID", req.UserBookID))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User book not found"})
		}
		log.Error("Failed to check UserBook existence", zap.Error(err), zap.Uint("userBookID", req.UserBookID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error checking user book"})
	}

	// TODO: Authorization check: Does the authenticated user own this UserBook?

	activity := models.ReadingActivity{
		UserBookID:  req.UserBookID,
		PagesRead:   req.PagesRead,
		StartPage:   req.StartPage,
		EndPage:     req.EndPage,
		Notes:       req.Notes,
		ReadingDate: req.ReadingDate,
	}

	// Use a transaction to ensure both activity creation and book update succeed or fail together.
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the reading activity
		if err := tx.Create(&activity).Error; err != nil {
			return err
		}

		// 2. Update the CurrentPage of the UserBook
		// We use the EndPage of the activity as the new CurrentPage of the book.
		// This assumes activities are logged chronologically.
		if err := tx.Model(&userBook).Update("current_page", activity.EndPage).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Error("Failed to create ReadingActivity and update UserBook", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create reading activity",
		})
	}

	log.Info("ReadingActivity created successfully", zap.Uint("activityID", activity.ID))
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reading activity created successfully",
		"data":    activity,
	})
}

// UpdateReadingActivity godoc
// @Summary Update a specific reading activity
// @Description Update details of a specific reading activity by its ID.
// @Tags ReadingActivity
// @Accept json
// @Produce json
// @Param activityId path int true "Reading Activity ID"
// @Param reading_activity_update body models.ReadingActivityUpdateRequest true "Reading Activity Update Payload"
// @Success 200 {object} fiber.Map{message=string, data=models.ReadingActivity}
// @Failure 400 {object} fiber.Map{message=string, errors=map[string]string}
// @Failure 404 {object} fiber.Map{message=string} "ReadingActivity not found"
// @Failure 500 {object} fiber.Map{message=string}
// @Router /reading-activities/{activityId} [put]
func (c *ReadingActivityController) UpdateReadingActivity(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	activityIDStr := ctx.Params("activityId")
	log.Info("ReadingActivityController.UpdateReadingActivity Begin", zap.String("activityID", activityIDStr))

	var req models.ReadingActivityUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body for ReadingActivity update", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request",
			"errors":  map[string]string{"body": "Failed to parse request body"},
		})
	}

	if err := c.Validate.Struct(&req); err != nil {
		log.Error("Validation failed for ReadingActivity update", zap.Error(err))
		validationErrors := make(map[string]string)
		for _, vErr := range err.(validator.ValidationErrors) {
			validationErrors[strings.ToLower(vErr.Field())] = vErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	var activity models.ReadingActivity
	if err := c.DB.First(&activity, activityIDStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("ReadingActivity not found for update", zap.String("activityID", activityIDStr))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Reading activity not found"})
		}
		log.Error("Failed to fetch ReadingActivity for update", zap.Error(err), zap.String("activityID", activityIDStr))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch reading activity"})
	}

	// TODO: Authorization check.

	// Apply updates from request
	if req.PagesRead != nil {
		activity.PagesRead = *req.PagesRead
	}
	if req.StartPage != nil {
		activity.StartPage = *req.StartPage
	}
	if req.EndPage != nil {
		activity.EndPage = *req.EndPage
	}
	if req.Notes != "" { // Assuming empty string means "not provided"
		activity.Notes = req.Notes
	}
	if !req.ReadingDate.IsZero() {
		activity.ReadingDate = req.ReadingDate
	}

	// Note: Updating an activity does not automatically update the UserBook's CurrentPage
	// as this can have complex side-effects (e.g., if this is not the latest activity).
	// This would require more complex business logic.

	if err := c.DB.Save(&activity).Error; err != nil {
		log.Error("Failed to update ReadingActivity in database", zap.Error(err), zap.Uint("activityID", activity.ID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update reading activity",
		})
	}

	log.Info("ReadingActivity updated successfully", zap.Uint("activityID", activity.ID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reading activity updated successfully",
		"data":    activity,
	})
}

// DeleteReadingActivity godoc
// @Summary Delete a specific reading activity
// @Description Permanently delete a specific reading activity by its ID.
// @Tags ReadingActivity
// @Accept json
// @Produce json
// @Param activityId path int true "Reading Activity ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 404 {object} fiber.Map{message=string} "ReadingActivity not found"
// @Failure 500 {object} fiber.Map{message=string}
// @Router /reading-activities/{activityId} [delete]
func (c *ReadingActivityController) DeleteReadingActivity(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	activityIDStr := ctx.Params("activityId")
	log.Info("ReadingActivityController.DeleteReadingActivity Begin", zap.String("activityID", activityIDStr))

	var activity models.ReadingActivity
	if err := c.DB.First(&activity, activityIDStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("ReadingActivity not found for deletion", zap.String("activityID", activityIDStr))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Reading activity not found"})
		}
		log.Error("Failed to fetch ReadingActivity for deletion", zap.Error(err), zap.String("activityID", activityIDStr))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch reading activity"})
	}

	// TODO: Authorization check.

	// Perform hard delete
	if err := c.DB.Unscoped().Delete(&activity).Error; err != nil {
		log.Error("Failed to delete ReadingActivity from database", zap.Error(err), zap.Uint("activityID", activity.ID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete reading activity",
		})
	}

	// Note: Deleting an activity does not automatically adjust the UserBook's CurrentPage.

	log.Info("ReadingActivity deleted successfully", zap.Uint("activityID", activity.ID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Reading activity deleted successfully"})
}

// GetAllReadingActivitiesForUserBook godoc
// @Summary Get all reading activities for a specific user book
// @Description Retrieve a list of all reading activities associated with a given UserBook ID.
// @Tags ReadingActivity
// @Accept json
// @Produce json
// @Param userBookId path int true "UserBook ID"
// @Success 200 {object} fiber.Map{message=string, data=[]models.ReadingActivity}
// @Failure 404 {object} fiber.Map{message=string} "UserBook not found"
// @Failure 500 {object} fiber.Map{message=string}
// @Router /userbooks/{userBookId}/activities [get]
func (c *ReadingActivityController) GetAllReadingActivitiesForUserBook(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	userBookIDStr := ctx.Params("userBookId")
	log.Info("ReadingActivityController.GetAllReadingActivitiesForUserBook Begin", zap.String("userBookID", userBookIDStr))

	// Validate UserBookID exists
	var userBook models.UserBook
	if err := c.DB.First(&userBook, userBookIDStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("UserBook not found when listing activities", zap.String("userBookID", userBookIDStr))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User book not found"})
		}
		log.Error("Failed to verify UserBook existence for listing activities", zap.Error(err), zap.String("userBookID", userBookIDStr))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error verifying user book"})
	}

	// TODO: Authorization check: Does the authenticated user own this UserBook?

	var activities []models.ReadingActivity
	if err := c.DB.Where("user_book_id = ?", userBook.ID).Order("reading_date DESC, created_at DESC").Find(&activities).Error; err != nil {
		log.Error("Failed to fetch reading activities from database", zap.Error(err), zap.Uint("userBookID", userBook.ID))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch reading activities",
		})
	}

	log.Info("Reading activities fetched successfully for UserBook", zap.Uint("userBookID", userBook.ID), zap.Int("count", len(activities)))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reading activities fetched successfully",
		"data":    activities,
	})
}

// GetReadingActivityByID godoc
// @Summary Get a specific reading activity by its ID
// @Description Retrieve details of a specific reading activity by its ID.
// @Tags ReadingActivity
// @Accept json
// @Produce json
// @Param activityId path int true "Reading Activity ID"
// @Success 200 {object} fiber.Map{message=string, data=models.ReadingActivity}
// @Failure 404 {object} fiber.Map{message=string} "ReadingActivity not found"
// @Failure 500 {object} fiber.Map{message=string}
// @Router /reading-activities/{activityId} [get]
func (c *ReadingActivityController) GetReadingActivityByID(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	activityIDStr := ctx.Params("activityId")
	log.Info("ReadingActivityController.GetReadingActivityByID Begin", zap.String("activityID", activityIDStr))

	var activity models.ReadingActivity
	// Preload UserBook to provide context.
	if err := c.DB.Preload("UserBook").First(&activity, activityIDStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("ReadingActivity not found by ID", zap.String("activityID", activityIDStr))
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Reading activity not found"})
		}
		log.Error("Failed to fetch ReadingActivity by ID from database", zap.Error(err), zap.String("activityID", activityIDStr))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch reading activity",
		})
	}

	// TODO: Authorization check: Does the authenticated user own the UserBook associated with this activity?
	// For example, after fetching activity:
	// var userBook models.UserBook
	// if err := c.DB.First(&userBook, activity.UserBookID).Error; err == nil {
	//   if authenticatedUserID != userBook.UserID && !IsAdmin(authenticatedUser) {
	//     return ctx.Status(fiber.StatusForbidden).JSON(...)
	//   }
	// } else { /* handle error fetching userbook for auth check */ }


	log.Info("ReadingActivity fetched successfully by ID", zap.Uint("activityID", activity.ID))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reading activity fetched successfully",
		"data":    activity,
	})
}

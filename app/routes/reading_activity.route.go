package routes

import (
	"ayo-baca-buku/app/controllers"
	// "ayo-baca-buku/app/middlewares" // Placeholder for auth middleware

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupReadingActivityRoutes(app *fiber.App, DB *gorm.DB) {
	// Instantiate the ReadingActivityController
	readingActivityController := controllers.NewReadingActivityController(DB)

	// Group for activities related to a specific user book
	// This route is for listing activities for a specific book
	userBookActivitiesRoutes := app.Group("/userbooks/:userBookId/activities")
	// Apply middleware here if needed, e.g., middlewares.AuthJWTMiddleware()
	userBookActivitiesRoutes.Get("/", readingActivityController.GetAllReadingActivitiesForUserBook)
	// Note: CreateReadingActivity currently expects UserBookID in the body.
	// A more RESTful approach for creation might be POST to this grouped route,
	// requiring controller adjustment to take UserBookID from path.
	// For now, Create will be a top-level route as per current controller design.

	// Group for general reading activity management (by activity ID)
	// Apply middleware here if needed
	activityRoutes := app.Group("/reading-activities")

	activityRoutes.Post("/", readingActivityController.CreateReadingActivity) // UserBookID in body
	activityRoutes.Get("/:activityId", readingActivityController.GetReadingActivityByID)
	activityRoutes.Put("/:activityId", readingActivityController.UpdateReadingActivity)
	activityRoutes.Delete("/:activityId", readingActivityController.DeleteReadingActivity)
}

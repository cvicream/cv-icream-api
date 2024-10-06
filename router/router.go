package router

import (
	"github.com/cvicream/cv-icream-api/handler"
	"github.com/cvicream/cv-icream-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	auth := api.Group("/auth")
	auth.Get("/user", middleware.Protected(), handler.GetCurrentUser)
	auth.Put("/user", middleware.Protected(), handler.UpdateCurrentUser)

	// Google auth
	googleAuth := api.Group("/auth/google")
	googleAuth.Get("/", handler.GoogleAuth)
	googleAuth.Get("/callback", handler.GoogleCallback)

	// LinkedIn auth
	linkedinAuth := api.Group("/auth/linkedin")
	linkedinAuth.Get("/", handler.LinkedInAuth)
	linkedinAuth.Get("/callback", handler.LinkedInCallback)

	// CV
	cv := api.Group("/cv", middleware.Protected())
	cv.Get("/", handler.GetAllCVs)
	cv.Get("/:id", handler.GetCV)
	cv.Post("/", middleware.Protected(), handler.CreateCV)
	cv.Put("/:id", middleware.Protected(), handler.UpdateCV)
	cv.Delete("/:id", middleware.Protected(), handler.DeleteCV)
}

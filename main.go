package main

import (
	"fmt"

	"github.com/cvicream/cv-icream-api/config"
	"github.com/cvicream/cv-icream-api/database"
	"github.com/cvicream/cv-icream-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	app.Use(cors.New())

	router.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	app.Listen(fmt.Sprintf(":%s", config.Config("PORT")))
	log.Info("Server running on port " + config.Config("PORT"))
}

package handler

import (
	"github.com/cvicream/cv-icream-api/model"
	"github.com/cvicream/cv-icream-api/service"

	"github.com/gofiber/fiber/v2"
)

func CreateSurvey(c *fiber.Ctx) error {
	survey := new(model.Survey)
	if err := c.BodyParser(survey); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid request body",
			})
	}

	if err := service.CreateSurvey(survey); err != nil {
			if err.Error() == "user not found" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"error": err.Error(),
					})
			}
			if err.Error() == "user already has a survey" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
						"error": err.Error(),
				})
		}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create survey",
			})
	}

	return c.JSON(survey)
}


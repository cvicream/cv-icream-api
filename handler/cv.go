package handler

import (
	"github.com/cvicream/cv-icream-api/model"
	"github.com/cvicream/cv-icream-api/service"

	"github.com/gofiber/fiber/v2"
)

func CreateCV(c *fiber.Ctx) error {
	cv := new(model.CV)
	if err := c.BodyParser(cv); err != nil {
		return c.SendStatus(400)
	}
	cv, err := service.CreateCV(*cv)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cv)
}

func GetAllCVs(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	cvs, err := service.GetCVs(userId)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cvs)
}

func GetCV(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(400)
	}

	cv, err := service.GetCV(userId, id)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cv)
}

func UpdateCV(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	cv := new(model.CV)
	if err := c.BodyParser(cv); err != nil {
		return c.SendStatus(400)
	}

	cv, err := service.UpdateCV(userId, cv)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cv)
}

func UpdateCVTitle(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	cv := new(model.CV)
	if err := c.BodyParser(cv); err != nil {
		return c.SendStatus(400)
	}

	cv, err := service.UpdateCVTitle(userId, cv)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cv)
}

func DeleteCV(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(400)
	}

	err := service.DeleteCV(userId, id)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

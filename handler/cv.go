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
	cvs, err := service.GetCVs()
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cvs)
}

func GetCV(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(400)
	}

	cv, err := service.GetCV(id)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cv)
}

func UpdateCV(c *fiber.Ctx) error {
	cv := new(model.CV)
	if err := c.BodyParser(cv); err != nil {
		return c.SendStatus(400)
	}

	cv, err := service.UpdateCV(cv)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(cv)
}

func DeleteCV(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(400)
	}

	err := service.DeleteCV(id)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

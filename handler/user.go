package handler

import (
	"github.com/cvicream/cv-icream-api/model"
	"github.com/cvicream/cv-icream-api/service"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	user, err := service.GetUserById(userId)
	if err != nil {
		return c.SendStatus(400)
	}
	return c.JSON(user)
}

func UpdateCurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(400)
	}

	user, err := service.UpdateUser(userId, user)
	if err != nil {
		return c.SendStatus(400)
	}
	return c.JSON(user)
}

func DeleteCurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	err := service.DeleteUser(userId)
	if err != nil {
		return c.SendStatus(400)
	}
	return nil
}

package handler

import (
	"github.com/cvicream/cv-icream-api/service"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUser(c *fiber.Ctx) error {
	jwtUser := c.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)

	user, err := service.GetUserById(userId)
	if err != nil {
		return c.SendStatus(400)
	}
	return c.JSON(user)
}

package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mnsh5/fiber-jwt-auth/src/database"
	"github.com/mnsh5/fiber-jwt-auth/src/models"
)

func User(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   fiber.Map{"user": user},
	})
}

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []models.User
	db.Find(&users)
	return c.JSON(users)
}

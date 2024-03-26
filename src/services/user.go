package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mnsh5/fiber-jwt-auth/src/database"
	"github.com/mnsh5/fiber-jwt-auth/src/models"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []models.User
	db.Find(&users)
	return c.JSON(users)
}

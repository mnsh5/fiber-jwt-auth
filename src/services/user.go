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

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(err)
	}

	db.Create(&user)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

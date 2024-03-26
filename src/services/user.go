package services

import (
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

package services

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mnsh5/fiber-jwt-auth/src/config"
	"github.com/mnsh5/fiber-jwt-auth/src/database"
	"github.com/mnsh5/fiber-jwt-auth/src/models"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c *fiber.Ctx) error {
	var payload *models.SignInInput
	db := database.DB

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var user models.User
	result := db.First(&user, "email = ?", strings.ToLower(payload.Email))

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "fail",
			"message": "The user with the email already exists",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
	}

	// Generar un nuevo par de claves Ed25519
	tokenByte := jwt.New(&jwt.SigningMethodEd25519{})
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	expDuration := time.Hour * 24

	claims["sub"] = user.ID
	claims["exp"] = now.Add(time.Hour * 24).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.Config("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
	}

	return nil
}

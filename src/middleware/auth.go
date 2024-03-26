package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mnsh5/fiber-jwt-auth/src/config"
	"github.com/mnsh5/fiber-jwt-auth/src/database"
	"github.com/mnsh5/fiber-jwt-auth/src/models"
)

func AuthMiddleware(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	}
	if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "You are not logged in"})
	}

	// Verificar la validez del token
	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		// Verificar el algoritmo de firma
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
		}
		// Devolver la clave secreta para validar el token
		return []byte(config.Config("SECRET_KEY")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Invalid token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token"})
	}

	var user models.User
	db := database.DB
	db.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	if float64(user.ID) != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "The user associated this token dosn't exists"})
	}

	c.Locals("user", &user)
	return c.Next()
}

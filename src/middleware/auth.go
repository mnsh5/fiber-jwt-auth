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

	// Si no hay token en el encabezado Authorization, intenta obtenerlo de la cookie
	if tokenString == "" {
		tokenString = c.Cookies("token")
	}

	// Si no hay token en el encabezado ni en la cookie, permitir que la solicitud continúe
	if tokenString == "" {
		return c.Next()
	}

	// Verificar la validez del token
	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		// Verificar el algoritmo de firma
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		// Devolver la clave secreta para validar el token
		return []byte(config.Config("SECRET_KEY")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token",
		})
	}

	// Obtener el usuario asociado con el token
	var user models.User
	db := database.DB
	db.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	// Verificar si el usuario asociado al token existe
	if user.ID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "The user associated with this token does not exist",
		})
	}

	// Agregar el usuario a los datos locales del contexto para que esté disponible en los controladores
	c.Locals("user", &user)
	return c.Next()
}

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mnsh5/fiber-jwt-auth/src/services"
)

func SetupRoutes(app *fiber.App) {
	// Grupo principal de la API
	api := app.Group("/api")

	// Subgrupo para la versión 1 de la API
	v1 := api.Group("/v1")

	// Rutas específicas de la versión 1 de la API
	v1.Post("/signup", services.SignUp)
	v1.Post("/signin", services.SignIn)
	v1.Post("/logout", services.Logout)
}

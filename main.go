package main

import (
	"github.com/gofiber/fiber/v2"

	controllers "root/controllers"
	initializers "root/initializers"
	models "root/models"
	middleware "root/middleware"
)

func main() {

	initializers.ConnectToDB()

	initializers.DB.MigrateTable(&models.User{})

	app := fiber.New()

	app.Post("/singup", controllers.SingUp)
	app.Post("/login", controllers.Login)
	app.Get("/validate",middleware.Auth, controllers.Validate)

	app.Listen(":3000")

}


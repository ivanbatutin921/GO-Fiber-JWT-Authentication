package main

import (
	"github.com/gofiber/fiber/v2"

	controllers "root/controllers"
	initializers "root/initializers"
	models "root/models"
	middleware "root/middleware"
)

// func initRoutes(app *fiber.App) {
// 	users := app.Group("/auth")
// 	{
// 		app.Post("/singup", controllers.SingUp)
// 		app.Post("/login", controllers.Login)
// 		app.Get("/validate", controllers.Validate)
// 		app.Get("/hello", controllers.Hello)

// 		authenticated := users.Group("/", middleware.AuthToken)
// 		{
// 			authenticated.Get("/hello", controllers.Hello)
// 		}
// 	}

// }

func main() {

	initializers.ConnectToDB()

	//initializers.DB.MigrateTable(&models.User{})
	initializers.DB.MigrateTable(&models.Token{})

	app := fiber.New()
	//initRoutes(app)

	app.Post("/singup", controllers.SingUp)
	app.Post("/login", controllers.Login)
	app.Get("/validate",middleware.AuthToken, controllers.Validate)
	app.Get("/hello",middleware.AuthRole,middleware.AuthToken, controllers.Hello)

	app.Listen(":3000")

}

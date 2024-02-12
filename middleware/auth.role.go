package middleware

import (
	"fmt"

	"root/initializers"
	"root/models"

	"github.com/gofiber/fiber/v2"
)

func AuthRole(c *fiber.Ctx) error {
	// body struct
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to read body",
		})
	}

	//find user	
	var user models.User
	initializers.DB.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status": "invalid email or password",
		})
	}
	
	role := string(body.Role)
	fmt.Print(body)
	if role == "admin" && role != "user" {
		c.Next()
	}else{
		return c.SendString("ХУЙ ДВА")
	}

	return nil
}

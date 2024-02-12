package controllers

import (
	//"errors"
	//"hash"
	"net/http"
	"os"
	"root/initializers"
	"root/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	

)

func SingUp(c *fiber.Ctx) error {
	// body struct
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	
	// parse body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to read body",
		})
	}

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to hash password",
		})
	}

	//create user
	user := models.User{Email: body.Email, Password: string(hash), Role: body.Role}
	result := initializers.DB.DB.Create(&user)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to create user",
		})
	}

	//respon
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status": "success",
	})
}

func Login(c *fiber.Ctx) error {
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

	//compare password and hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "invalid email or password",
		})
	}

	//generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to generate token",
		})
	}

	//send token
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 12),
	})

	

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status": tokenString,
		"message": user.Role,
	})
}

func Validate(c *fiber.Ctx) error {
	c.JSON(fiber.Map{
		"message": "I entered",
	})
	return nil
}

func Hello (c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
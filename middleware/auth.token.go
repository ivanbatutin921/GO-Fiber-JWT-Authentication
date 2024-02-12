package middleware

import (
	"fmt"
	"log"
	"os"
	"root/initializers"
	"root/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthToken(c *fiber.Ctx) error {
	// get token from cookie
	tokenString := c.Cookies("jwt")
	if tokenString == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRETKEY")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check expiration
		if float64(time.Now().Unix()) > float64(claims["exp"].(float64)) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		// find user
		var user models.User
		initializers.DB.DB.Find(&user, claims["sub"])
		if user.ID == 0 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// add user to context
		c.Locals("user", user)

		// move to next middleware
		c.Next()

	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return nil
}

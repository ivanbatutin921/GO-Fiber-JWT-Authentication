package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
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

	//refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to generate refresh token",
		})
	}

	// save refresh token to database
	refreshTokenRecord := models.Token{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Expiry:       time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	result = initializers.DB.DB.Create(&refreshTokenRecord)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to save refresh token",
		})
	}

	//add refresh token to cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 12),
	})

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
		"exp": time.Now().Add(time.Minute * 15).Unix(),
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

	//generate refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to generate refresh token",
		})
	}
	
	//update access token to database
	updateResult := initializers.DB.DB.Model(&models.Token{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Expiry:       time.Now().Add(time.Hour * 24 * 30).Unix(),
	}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
		"refresh_token": refreshToken,
		"expiry":        time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	if updateResult.Error != nil {
		log.Println(updateResult.Error)
		return c.Status(500).JSON(fiber.Map{
			"status": "fail to update refresh token",
		})
	}
	//respon
	return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":    "success",
		"access_token": tokenString,
	})	

}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32) // generate 32 random bytes

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(b), nil

}

func Validate(c *fiber.Ctx) error {
	c.JSON(fiber.Map{
		"message": "I entered",
	})
	return nil
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// // проверьте, существует ли токен обновления и является ли он действительным
// var refreshTokenRecord models.Token
// initializers.DB.DB.Where("refresh_token = ? AND expiry > ?", c.Cookies("refresh_token"), time.Now().Unix()).First(&refreshTokenRecord)
// if refreshTokenRecord.ID != 0 {
// 	// проверьте токен обновления
// 	token, err := jwt.Parse(c.Cookies("refresh_token"), func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC);
// 			!ok {
// 				return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
// 			}
// 		return []byte(os.Getenv("SECRETKEY")), nil
// 	})
// 	if err != nil {
// 		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
// 			"status": "refresh token is expired or invalid",
// 		})
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// find the user
// 		var user models.User
// 		initializers.DB.DB.First(&user, "id = ?", int(claims["sub"].(float64)))
// 		if user.ID != 0 {
// 			// generate new access token
// 			token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 				"sub": user.ID,
// 				"exp": time.Now().Add(time.Minute * 15).Unix(),
// 			})
// 			tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
// 			if err != nil {
// 				return c.Status(500).JSON(fiber.Map{
// 					"status": "fail to generate token",
// 				})
// 			}
// 			// send token
// 			c.Cookie(&fiber.Cookie{
// 				Name:     "jwt",
// 				Value:    tokenString,
// 				HTTPOnly: true,
// 				Expires:  time.Now().Add(time.Hour * 12),
// 			})
// 			return c.Status(http.StatusOK).JSON(fiber.Map{
// 				"status": tokenString,
// 				"message": user.Role,
// 			})
// 		}
// 	}
// }
// return c.Next()

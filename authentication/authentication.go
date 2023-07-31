package authentication

import (

	//inbuilt package
	"net/http"
	"time"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// setting authorization for the users
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).SendString("Missing token")

		}
		for index, char := range tokenString {
			if char == ' ' {
				tokenString = tokenString[index+1:]
			}
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).SendString("Invalid token")
		}

		// Check whether token is expired or not
		check, ok := claims["exp"].(int64)
		if ok && check < time.Now().Unix() {
			return c.Status(http.StatusUnauthorized).SendString("Expired token")
		}
		return c.Next()
	}
}

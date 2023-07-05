package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("unauthorized")
	}

	if err := parseToken(token); err != nil {
		return err
	}

	return c.Next()
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["email"] != "" &&
		claims["id"] != "" {
		valid, err := time.Parse(time.RFC3339Nano, claims["validTill"].(string))
		// "2023-07-06T00:44:05.818238369+03:00"
		if err != nil {
			fmt.Println("can't parse time", err)
		} else if valid.Before(time.Now()) {
			fmt.Println("time of the token is wrong")
		} else {
			return nil
		}
	}
	return fmt.Errorf("unauthorized")
}

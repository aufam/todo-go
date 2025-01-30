package models

import (
	"os"
	"time"
	"todo-go/core"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

type JWTResponse struct {
	Token string `json:"token"`
}

func JWTEncode(userID string) (JWTResponse, error) {
	secretKey, err := core.LoadEnv("JWT_SECRET_KEY")
	if err != nil {
		return JWTResponse{}, err
	}

	expStr, err := core.LoadEnv("JWT_EXPIRATION_LENGTH")
	if err != nil {
		return JWTResponse{}, err
	}

	exp, err := time.ParseDuration(expStr)
	if err != nil {
		return JWTResponse{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	})

	res, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return JWTResponse{}, err
	}

	return JWTResponse{Token: res}, nil
}

func JWTMiddleware(c *fiber.Ctx) error {
	middleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		ContextKey: "jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "Invalid or expired token"})
		},
	})

	if err := middleware(c); err != nil {
		return err
	}

	userToken, ok := c.Locals("jwt").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "Invalid token"})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "Invalid token structure"})
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().After(time.Unix(int64(exp), 0)) {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "Token expired"})
		}
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "Invalid user ID"})
	}

	c.Locals("userID", userID)
	return c.Next()
}

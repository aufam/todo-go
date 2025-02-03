package models

import (
	"fmt"
	"os"
	"strings"
	"time"
	"todo-go/core"

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
	checkToken := func(auth string) string {
		authentication := c.Get(auth)
		if strings.HasPrefix(authentication, "Bearer ") {
			return strings.TrimPrefix(authentication, "Bearer ")
		}
		return ""
	}

	tokenString := checkToken("Authentication")
	if tokenString == "" {
		tokenString = checkToken("authentication")
	}
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "You are not logged in"})
	}

	claims := JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(MakeErrorResponse(err))
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Err: "Token expired"})
	}

	c.Locals("userID", claims.UserID)
	return c.Next()
}

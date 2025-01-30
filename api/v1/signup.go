package v1

import (
	"context"
	"todo-go/db"
	"todo-go/models"

	"github.com/gofiber/fiber/v2"
)

func Signup(ctx context.Context, database db.Database, user models.UserForm) (status int, response any) {
	userID, err := database.AddUser(ctx, user)
	if err != nil {
		return fiber.StatusInternalServerError, models.MakeErrorResponse(err)
	}

	token, err := models.JWTEncode(userID)
	if err != nil {
		database.DelUser(ctx, userID)
		return fiber.StatusInternalServerError, models.MakeErrorResponse(err)
	}

	return fiber.StatusOK, token
}

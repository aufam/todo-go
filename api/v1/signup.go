package v1

import (
	"todo-go/db"
	"todo-go/models"

	"github.com/gofiber/fiber/v2"
)

func Signup(database db.Database, user models.UserForm) (status int, response any) {
	userID, err := database.AddUser(user)
	if err != nil {
		return fiber.StatusInternalServerError, models.MakeErrorResponse(err)
	}

	token, err := models.JWTEncode(userID)
	if err != nil {
		return fiber.StatusInternalServerError, models.MakeErrorResponse(err)
	}

	return fiber.StatusOK, token
}

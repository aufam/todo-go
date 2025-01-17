package v1

import (
	"todo-go/db"
	"todo-go/models"

	"github.com/gofiber/fiber/v2"
)

func Login(database db.Database, user models.UserForm) (status int, response any) {
	userID, err := database.GetUserID(user)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	token, err := models.JWTEncode(userID)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, token
}

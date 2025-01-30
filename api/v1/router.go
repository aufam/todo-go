package v1

import (
	"todo-go/db"
	"todo-go/models"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, database db.Database) {
	api := app.Group("/api/v1")

	api.Use(func(c *fiber.Ctx) error {
		c.Locals("db", database)
		return c.Next()
	})

	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("Todo App v1")
	})

	api.Post("/user/signup", func(c *fiber.Ctx) error {
		var user models.UserForm
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.MakeErrorResponse(err))
		}

		status, response := Signup(c.Context(), c.Locals("db").(db.Database), user)
		return c.Status(status).JSON(response)
	})

	api.Post("/user/login", func(c *fiber.Ctx) error {
		var user models.UserForm
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.MakeErrorResponse(err))
		}

		status, response := Login(c.Context(), c.Locals("db").(db.Database), user)
		return c.Status(status).JSON(response)
	})

	api.Get("/user", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)

		user, err := database.GetUser(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(user)
	})

	api.Delete("/user", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)

		err := database.DelUser(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/todos", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)

		todos, err := database.GetTodos(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(todos)
	})

	api.Post("/todo", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)

		var todo models.TodoForm
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.MakeErrorResponse(err))
		}

		todoID, err := database.AddTodo(c.Context(), userID, todo)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(todoID)
	})

	api.Get("/todo/:id", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)
		todoID := c.Params("id")

		todo, err := database.GetTodo(c.Context(), userID, todoID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(todo)
	})

	api.Put("/todo/:id", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)
		todoID := c.Params("id")

		var todo models.TodoForm
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.MakeErrorResponse(err))
		}

		err := database.ModTodo(c.Context(), userID, todoID, todo)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Delete("/todo/:id", models.JWTMiddleware, func(c *fiber.Ctx) error {
		database := c.Locals("db").(db.Database)
		userID := c.Locals("userID").(string)
		todoID := c.Params("id")

		err := database.DelTodo(c.Context(), userID, todoID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeErrorResponse(err))
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

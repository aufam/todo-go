package main

import (
	"fmt"
	v1 "todo-go/api/v1"
	"todo-go/core"
	"todo-go/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	database, err := db.OpenDefault()
	if err != nil {
		panic(err)
	}

	defer database.Close()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	version, err := core.GetVersion()
	if err != nil {
		panic(err)
	}

	if version == "v1" {
		v1.SetupRoutes(app, database)
	} else {
		panic(fmt.Errorf("Unknown version `%s`", version))
	}

	app.Listen(":8000")
}

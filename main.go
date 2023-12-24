package main

import (
	"github.com/Epic55/go_project_jwt_oauth/db"
	"github.com/Epic55/go_project_jwt_oauth/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.DBconn()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true, //Very important while using a HTTPonly Cookie, frontend can easily get and return back the cookie.
	}))
	routes.Setup(app)
	app.Listen("localhost:8080")
}

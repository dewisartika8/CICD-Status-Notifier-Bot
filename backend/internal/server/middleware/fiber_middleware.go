/*
File: Request Middleware
@author -riff-
Date : 27-02-2025
*/
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// FiberMiddleware provide Fiber's built-in middlewares.
func FiberMiddleware(a *fiber.App) {
	a.Use(
		recover.New(recover.Config{EnableStackTrace: true}),
		// Add CORS to each route.
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
			AllowMethods: "GET,POST,PUT,HEAD",
		}),
		// Add helmet middleware
		helmet.New(),
	)
}

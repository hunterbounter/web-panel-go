package acl

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func Authorized() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// Admin ACL
func AdminAuthorized() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.Next()

	}
}

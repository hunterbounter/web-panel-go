package acl

import (
	"github.com/gofiber/fiber/v2"
	hunterbounter_session "hunterbounter.com/web-panel/pkg/session"
	"log"
)

func Unauthorized() fiber.Handler {
	return func(c *fiber.Ctx) error {

		//print session
		sess, err := hunterbounter_session.HunterSession.Get(c)
		log.Println("Disallow Anon -> Session: ", sess)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if sess.Get("userId") != nil {
			return c.Redirect("/dashboard")
		}

		return c.Next()
	}
}

func Authorized() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := hunterbounter_session.HunterSession.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if sess.Get("userId") == nil {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}

// Admin ACL
func AdminAuthorized() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.Next()

	}
}

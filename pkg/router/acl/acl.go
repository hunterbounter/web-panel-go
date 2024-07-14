package acl

import (
	"github.com/gofiber/fiber/v2"
	hunterbounter_session "hunterbounter.com/web-panel/pkg/session"
	"log"
)

func Unauthorized() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var sessionFound bool
		//print session
		sess, err := hunterbounter_session.GetSessionValue(c, hunterbounter_session.SESSION_KEY)
		if err != nil {
			sessionFound = false
		}

		if sess != nil && sess.UserID != "" {
			sessionFound = true
		}

		if sessionFound {
			return c.Redirect("/dashboard")
		}
		return c.Next()
	}
}

func Authorized() fiber.Handler {
	return func(c *fiber.Ctx) error {

		sess, err := hunterbounter_session.GetSessionValue(c, hunterbounter_session.SESSION_KEY)
		if err != nil {
			log.Println("Session Error: ", err)
			return c.Redirect("/login")
		}

		if sess.UserID == "" {
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

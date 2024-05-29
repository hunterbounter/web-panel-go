package controller

import "github.com/gofiber/fiber/v2"

func DashboardGET(c *fiber.Ctx) error {

	return RenderTemplate(c, "panel/dashboard", fiber.Map{
		"Title": "Dashboard",
	})

}

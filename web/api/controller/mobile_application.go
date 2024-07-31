package controller

import "github.com/gofiber/fiber/v2"

func MobileScanGET(c *fiber.Ctx) error {
	return RenderTemplate(c, "panel/mobile_application", fiber.Map{
		"Title": "Mobile Application Scan",
	})
}

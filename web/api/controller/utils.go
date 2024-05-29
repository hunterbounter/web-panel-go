package controller

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func CheckIsDev() bool {
	if os.Getenv("USER") == "eminsargin" {
		return true
	}
	return false
}

func RenderTemplate(c *fiber.Ctx, name string, data fiber.Map) error {

	var SiteURI string
	if CheckIsDev() {
		SiteURI = "http://localhost:9000/"
	} else {
		SiteURI = "https://panel.hunterbounter.com/"
	}

	// Global variables
	globalData := fiber.Map{
		"SiteURI": SiteURI,
	}

	// Add global data to the existing data map
	for k, v := range globalData {
		if _, exists := data[k]; !exists {
			data[k] = v
		}
	}

	// Custom render process
	return c.Render(name, data)
}

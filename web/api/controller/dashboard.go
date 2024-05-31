package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"log"
)

const BACKEND_URL = "http://localhost:9001/"

func DashboardGET(c *fiber.Ctx) error {

	return RenderTemplate(c, "panel/dashboard", fiber.Map{
		"Title": "Dashboard",
	})

}

func SaveTarget(c *fiber.Ctx) error {

	log.Println("Coming from SaveTarget")

	// proxy request to backend

	return proxy.Forward(BACKEND_URL + "api/targets")(c)

}

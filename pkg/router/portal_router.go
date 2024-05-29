package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
)

type HttpRouter struct {
}

func (h HttpRouter) InstallRouter(app *fiber.App) {

	api := app.Group("", logger.New())

	// Check if the server is up
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Pong", nil))
	})

}

func NewHunterBounterApiRouter() *HttpRouter {
	return &HttpRouter{}
}

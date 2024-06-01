package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/router/acl"
	"hunterbounter.com/web-panel/web/api/controller"
)

type HttpRouter struct {
}

func (h HttpRouter) InstallRouter(app *fiber.App) {

	portal := app.Group("/", logger.New())

	api := app.Group("/api", logger.New())

	telemetry := app.Group("/telemetry", logger.New())

	portal.Get("/", acl.Unauthorized(), controller.DashboardGET)

	portal.Get("/vulnerability-report", acl.Unauthorized(), controller.VulnerabilityReportGET)

	/*
		Api Endpoints
	*/

	api.Post("/target/save", acl.Unauthorized(), controller.SaveTarget)

	/*
		Telemetry Endpoints
	*/

	telemetry.Get("/save", acl.Unauthorized(), controller.SaveTelemetryGet)
	telemetry.Post("/save", acl.Unauthorized(), controller.SaveTelemetry)

	/*
		Get Scanned Targets
	*/

	portal.Post("/target", acl.Unauthorized(), controller.GetTargets)

	// Check if the server is up
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Pong", nil))
	})

}

func NewHunterBounterApiRouter() *HttpRouter {
	return &HttpRouter{}
}

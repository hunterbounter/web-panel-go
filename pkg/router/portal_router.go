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

	portal.Get("/login", acl.Unauthorized(), controller.LoginGet)
	api.Post("/login", acl.Unauthorized(), controller.LoginPost)
	portal.Get("/", acl.Authorized(), controller.DashboardGET)

	portal.Get("/vulnerability-report", acl.Authorized(), controller.VulnerabilityReportGET)
	portal.Get("/vulnerability-report/:id", acl.Authorized(), controller.VulnerabilityReportDetailGET)

	api.Post("/scan/start", acl.Authorized(), controller.StartScan)

	/*
		Api Endpoints
	*/

	api.Post("/target/save", acl.Authorized(), controller.SaveTarget)

	/*
		Telemetry Endpoints
	*/

	telemetry.Post("/save", acl.Unauthorized(), controller.SaveTelemetry)

	/*
		Get Scanned Targets
	*/

	portal.Post("/target", acl.Unauthorized(), controller.GetTargets)

	/*
		Scan Results
	*/
	portal.Post("/scan_results/save", acl.Unauthorized(), controller.ScanResultPOST)

	api.Post("/agent/kill", acl.Authorized(), controller.KillAgent)

	// Check if the server is up
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Pong", nil))
	})

}

func NewHunterBounterApiRouter() *HttpRouter {
	return &HttpRouter{}
}

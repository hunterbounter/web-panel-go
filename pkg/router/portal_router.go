package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/router/acl"
	"hunterbounter.com/web-panel/web/api/controller"
	"time"
)

type HttpRouter struct {
}

func (h HttpRouter) InstallRouter(app *fiber.App) {

	rateLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("Are u kidding me?")
		},
	})

	app.Use(rateLimiter)
	portal := app.Group("/", logger.New())

	api := app.Group("/api", logger.New())

	telemetry := app.Group("/telemetry", logger.New())

	portal.Get("/login", acl.Unauthorized(), controller.LoginGet)
	portal.Get("/admin", acl.Authorized(), func(ctx *fiber.Ctx) error {
		//redirect https://www.youtube.com/watch?v=M5VXCixTdEg
		return ctx.Redirect("https://www.youtube.com/watch?v=M5VXCixTdEg")
	})
	api.Post("/login", acl.Unauthorized(), controller.LoginPost)
	portal.Get("/", acl.Authorized(), controller.DashboardGET)

	portal.Get("/targets", acl.Authorized(), controller.GetTargetsView)
	portal.Get("/vulnerability-report", acl.Authorized(), controller.VulnerabilityReportGET)

	portal.Get("/vulnerability-report/zap/:id", acl.Authorized(), controller.ZapReportDetailGET)
	portal.Get("/vulnerability-report/openvas/:id", acl.Authorized(), controller.OpenVasReportDetailGET)
	portal.Get("/vulnerability-report/nuclei/:id", acl.Authorized(), controller.NucleiReportDetailGET)

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

	portal.Post("/scan_results/save_pdf", acl.Unauthorized(), controller.ScanResultSavePDF)

	portal.Post("/scan_results/openvas/save", acl.Unauthorized(), controller.ScanResultOpenVASPOST)

	portal.Post("/scan_results/nuclei/save", acl.Unauthorized(), controller.ScanResultNucleiPOST)

	api.Post("/agent/kill", acl.Authorized(), controller.KillAgent)

	/*
		Mobile Application Scan
	*/

	portal.Get("/mobile-scan", acl.Authorized(), controller.MobileScanGET)

	/*
		Upload File
	*/
	portal.Post("/mobile-scan/upload", acl.Authorized(), controller.UploadMobileAppFile)

	/*
		Leak Data Search
	*/
	portal.Get("/leak-data", acl.Authorized(), controller.LeakDataSearchGET)
	portal.Post("/leak-data", acl.Authorized(), controller.LeakDataSearchPost)

	// Check if the server is up
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Pong", nil))
	})

}

func NewHunterBounterApiRouter() *HttpRouter {
	return &HttpRouter{}
}

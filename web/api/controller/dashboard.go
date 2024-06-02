package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/utils"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

const BACKEND_URL = "http://localhost:9001/"

func DashboardGET(c *fiber.Ctx) error {

	var totalReportCount = model.GetTotalReportCount()

	var monitorData = model.GetMonitorData()

	var totalWaitingScanCount = model.GetWaitingTargetDomainCount()

	var totalScannedTargetCount = model.GetTotalScannedTargetCount()

	return RenderTemplate(c, "panel/dashboard", fiber.Map{
		"Title":                   "Dashboard",
		"MonitorData":             monitorData,
		"TotalReportCount":        totalReportCount,
		"TotalWaitingTargetCount": totalWaitingScanCount,
		"TotalScannedTargetCount": totalScannedTargetCount,
	})

}

func SaveTarget(c *fiber.Ctx) error {

	log.Println("Coming from SaveTarget")

	targetsRaw := c.FormValue("targets")
	if targetsRaw == "" {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Targets cannot be empty", nil))
	}

	// targetsRaw'ı JSON listesine dönüştür
	var targetsList []string
	if err := json.Unmarshal([]byte(targetsRaw), &targetsList); err != nil {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Targets cannot be empty", nil))
	}

	for _, target := range targetsList {
		// Save target to database
		isValidDomain := utils.IsValidDomain(target)
		if isValidDomain {
			model.SaveTarget(target, 1, 0) // Domain And Added
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Targets saved successfully", nil))
	//return proxy.Forward(BACKEND_URL + "api/targets")(c)

}

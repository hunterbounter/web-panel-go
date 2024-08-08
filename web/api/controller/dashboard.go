package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/utils"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

func DashboardGET(c *fiber.Ctx) error {

	var totalReportCount = model.GetTotalReportCount()

	var monitorData = model.GetMonitorData()

	var totalWaitingScanCount int
	var zapWaitingDomainCount = model.GetWaitingZAPTargetDomainCount()

	var openvasWaitingDomainCount = model.GetWaitingOpenVASTargetDomainCount()

	var nucleiWaitingDomainCount = model.GetWaitingNucleiTargetDomainCount()

	totalWaitingScanCount = zapWaitingDomainCount + openvasWaitingDomainCount + nucleiWaitingDomainCount

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

	selectedAgents := c.FormValue("selectedAgents")

	log.Println("Selected Agents : ", selectedAgents)

	if selectedAgents == "" {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Please select an agent", nil))
	}

	var selectedAgentList []string
	if err := json.Unmarshal([]byte(selectedAgents), &selectedAgentList); err != nil {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Please select an agent", nil))
	}

	if len(selectedAgentList) == 0 {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Please select an agent", nil))
	}

	var isZap, isOpenVas, isNuclei, isMobSF bool
	for _, agent := range selectedAgentList {
		if agent == "" {
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Please select an agent", nil))
		}
		if agent != "OpenVas" && agent != "ZapProxy" && agent != "Nuclei" {
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Please select an agent", nil))
		}
		if agent == "OpenVas" {
			isOpenVas = true
		}
		if agent == "ZapProxy" {
			isZap = true
		}
		if agent == "Nuclei" {
			isNuclei = true
		}
		if agent == "MobSF" {
			isMobSF = true
		}

	}

	log.Println("isZap : ", isZap)
	log.Println("isOpenVas : ", isOpenVas)
	log.Println("isNuclei : ", isNuclei)
	log.Println("isMobSF : ", isMobSF)

	targetsRaw := c.FormValue("targets")
	if targetsRaw == "" {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Targets cannot be empty", nil))
	}

	var targetsList []string
	if err := json.Unmarshal([]byte(targetsRaw), &targetsList); err != nil {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Targets cannot be empty", nil))
	}

	for _, target := range targetsList {
		// Save target to database
		log.Println("Target : ", target)
		if isZap {

			// check is start http or https
			if target[:4] == "http" || target[:5] == "https" {
				model.SaveTarget(target, 1, 0) // Domain And Added
			}

			isValidDomain := utils.IsValidDomain(target)
			if isValidDomain {
				model.SaveTarget(target, 1, 0) // Domain And Added
			}

		}
		if isOpenVas {
			isValidIP := utils.IsValidIP(target)
			if isValidIP {
				model.SaveTarget(target, 2, 0) // IP And Added
			}
		}
		if isNuclei {
			var isValid bool
			if utils.IsValidDomainNormal(target) || utils.IsValidIP(target) || utils.IsValidDomain(target) {
				isValid = true
			}
			if isValid {
				model.SaveTarget(target, 3, 0) // Nuclei And Added
			}
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Targets saved successfully", nil))
	//return proxy.Forward(BACKEND_URL + "api/targets")(c)

}

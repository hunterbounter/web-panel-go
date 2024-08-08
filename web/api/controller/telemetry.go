package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"log"
)

func CheckTelemetryIsExists(hostname string) bool {

	var whereCondition = map[string]interface{}{"hostname": hostname}

	// Check if the record is exist
	dbRecords, err := database.Select("machine_monitor", whereCondition)

	if err != nil {
		log.Println("Error while checking record is exist", err)
		return false
	}

	if len(dbRecords) > 0 {
		return true
	}
	return false
}

func SaveTelemetry(c *fiber.Ctx) error {
	log.Println("Telemetry Income")

	// Get Request JSON
	body := c.Body()

	var telemetry map[string]interface{}
	err := json.Unmarshal(body, &telemetry)
	if err != nil {

		log.Println("Error Telemertry Data", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error parsing request", nil))
	}

	log.Println("Telemetry data -> ", hunterbounter_json.ToString(telemetry))

	// Print telemetry data as map

	var dbData = map[string]interface{}{
		"hostname":  telemetry["hostname"],
		"last_seen": database.CurrentTime(),
	}
	switch telemetry["telemetry_type"] {
	case "zap":
		dbData["type"] = 1 // 1 Zap 2 OpenVas 3 VPN
		dbData["active_scan_count"] = telemetry["active_scan_count"]
		dbData["service_status"] = telemetry["zap_status"]
	case "openvas":
		dbData["type"] = 2 // 1 Zap 2 OpenVas 3 VPN
		dbData["active_scan_count"] = telemetry["active_scan_count"]
		dbData["service_status"] = telemetry["openvas_status"]
	case "nuclei":
		dbData["type"] = 3 // 1 Zap 2 OpenVas 3 nuclei
		dbData["active_scan_count"] = telemetry["active_scan_count"]
		dbData["service_status"] = telemetry["nuclei_status"]
	case "mobsf":
		dbData["type"] = 4 // 1 Zap 2 OpenVas 3 nuclei
		dbData["active_scan_count"] = telemetry["active_scan_count"]
		dbData["service_status"] = telemetry["mobsf_status"]
	}

	if CheckTelemetryIsExists(telemetry["hostname"].(string)) {

		// Update telemetry data
		_, err := database.Update("machine_monitor", dbData, map[string]interface{}{"hostname": telemetry["hostname"]})
		if err != nil {
			log.Println("Telemetry data -> ", telemetry)
			log.Println("Error while updating telemetry data", err)

		}
	} else {
		log.Println("EKleniyorrrr")
		// Insert telemetry data
		_, err := database.Insert("machine_monitor", dbData, false)
		log.Println("Added telemetry data -> ", telemetry)
		if err != nil {
			log.Println("Telemetry data -> ", telemetry)
			log.Println("Error while inserting telemetry data", err)
		}
	}

	return c.SendString("Telemetry data received successfully")
}

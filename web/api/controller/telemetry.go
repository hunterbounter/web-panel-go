package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	"log"
)

func SaveTelemetryGet(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"succes": "true",
	})
}
func SaveTelemetry(c *fiber.Ctx) error {
	log.Println("Telemetry Income")

	// Get Request JSON
	body := c.Body()

	var telemetry map[string]interface{}
	err := json.Unmarshal(body, &telemetry)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Print telemetry data as map
	log.Printf("Received telemetry data: %+v\n", telemetry)

	log.Println(hunterbounter_json.ToString(telemetry))

	return c.SendString("Telemetry data received successfully")
}

package controller

import (
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/docker"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

type KillAgentRequest struct {
	Hostname string `json:"hostname" form:"hostname"`
}

func KillAgent(c *fiber.Ctx) error {
	var request KillAgentRequest

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Invalid Request", nil))
	}

	dockerManager := docker.NewDockerManager()

	err = dockerManager.KillContainer(request.Hostname)
	if err != nil {
		log.Println("REMOVE CONTAINER ERROR", err)
		return err
	}

	model.RemoveMonitorData(request.Hostname)
	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Agent Killed", nil))
}

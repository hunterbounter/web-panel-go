package controller

import (
	"github.com/gofiber/fiber/v2"
	hunterbounter_docker "hunterbounter.com/web-panel/pkg/docker"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

const DOCKER_DOMAIN_IMAGE_NAME = "hunter_bounter_zapv1"

const MAX_DOCKER_CONTAINER_COUNT = 5

func StartScan(c *fiber.Ctx) error {

	// First find total waiting scan count

	var totalWaitingScanCount = model.GetWaitingTargetDomainCount()

	dockerManager := hunterbounter_docker.NewDockerManager()

	containers, err := dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME) // List all running containers

	log.Println(containers)
	log.Println("Total running containers", len(containers))

	// Eğer çalışan container sayısı max container sayısından küçükse eşitlenene kadar yeni container başlat
	if len(containers) < MAX_DOCKER_CONTAINER_COUNT {
		for i := 0; i < MAX_DOCKER_CONTAINER_COUNT-len(containers); i++ {
			imageName := "hunter_bounter_zapv1"
			user := "root"
			port := "5002:5002"
			dns := "1.1.1.1"

			containerID, err := dockerManager.RunContainer(imageName, user, port, dns)
			log.Println("Running container ID", containerID)
			if err != nil {
				log.Println("Error while starting container", err)
				return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while starting container", nil))
			}
		}
	}

	lastContainersCount, err := dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME) // List all running containers

	log.Println("last container count", len(lastContainersCount))

	if err != nil {
		log.Println("Error while listing containers", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while listing containers", nil))
	}
	log.Println("Total running containers", len(containers))

	// Type 1 Domain
	// Status 0 Added
	model.UpdateTargetsStatus(0, 1, 1)

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan started", totalWaitingScanCount))
}

package controller

import (
	"github.com/gofiber/fiber/v2"
	hunterbounter_docker "hunterbounter.com/web-panel/pkg/docker"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

const DOCKER_DOMAIN_IMAGE_NAME = "hunter_bounter_zapv1"

const maxContainers = 5
const targetsPerContainer = 10

func StartScan(c *fiber.Ctx) error {

	// First find total waiting scan count

	var zapWaitingScanCount = model.GetWaitingZAPTargetDomainCount()
	log.Println("ZAP Waiting Scan Count", zapWaitingScanCount)
	var openvasWaitingScanCount = model.GetWaitingOpenVASTargetDomainCount()
	log.Println("OpenVAS Waiting Scan Count", openvasWaitingScanCount)

	var totalWaitingScanCount = zapWaitingScanCount + openvasWaitingScanCount

	requiredContainers := (totalWaitingScanCount + targetsPerContainer - 1) / targetsPerContainer // Round up to the nearest whole number

	if requiredContainers > maxContainers {
		requiredContainers = maxContainers
	}

	dockerManager := hunterbounter_docker.NewDockerManager()

	containers, err := dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME) // List all running containers

	if err != nil {
		log.Println("Error while listing running containers", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while listing running containers", nil))
	}

	// If the number of running containers is less than the required number of containers, start new containers until it is equalized
	log.Println("Total running containers", len(containers))
	log.Println("Required container count", requiredContainers)
	if len(containers) < requiredContainers {
		zapContainersToStart := (zapWaitingScanCount + targetsPerContainer - 1) / targetsPerContainer
		log.Println("zapContainersToStart", zapContainersToStart)
		openvasContainersToStart := (openvasWaitingScanCount + targetsPerContainer - 1) / targetsPerContainer
		log.Println("openvasContainersToStart", openvasContainersToStart)
		for i := 0; i < zapContainersToStart; i++ {
			imageName := "hunter_bounter_zapv1"
			user := "root"
			port := "5002:5002"
			dns := "1.1.1.1"

			containerID, err := dockerManager.RunContainer(imageName, user, port, dns)
			log.Println("Running ZAP container ID", containerID)
			if err != nil {
				log.Println("Error while starting ZAP container", err)
				return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while starting ZAP container", nil))
			}
		}

		for i := 0; i < openvasContainersToStart; i++ {
			imageName := "hunterbounter:openvasv1"

			containerID, err := dockerManager.RunContainerNoCommand(imageName)
			log.Println("Running OpenVAS container ID", containerID)
			if err != nil {
				log.Println("Error while starting OpenVAS container", err)
				return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while starting OpenVAS container", nil))
			}
		}
	}

	lastContainersCount, err := dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME) // List all running containers
	if err != nil {
		log.Println("Error while listing running containers", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while listing running containers", nil))
	}

	log.Println("last container count", len(lastContainersCount))

	if err != nil {
		log.Println("Error while listing containers", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while listing containers", nil))
	}
	log.Println("Total running containers", len(containers))

	// Type 1 Domain
	// Status 0 Added
	model.UpdateTargetsStatus(0, 1, 1)
	model.UpdateTargetsStatus(0, 2, 1)

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan started", totalWaitingScanCount))
}

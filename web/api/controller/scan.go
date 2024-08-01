package controller

import (
	"github.com/gofiber/fiber/v2"
	hunterbounter_docker "hunterbounter.com/web-panel/pkg/docker"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

const DOCKER_DOMAIN_IMAGE_NAME_ZAP = "hunter_bounter_zapv1"
const DOCKER_DOMAIN_IMAGE_NAME_OPENVAS = "hunterbounter:openvasv1"
const DOCKER_DOMAIN_IMAGE_NAME_NUCLEI = "hunterbounter:nucleiv1"

const maxContainers = 5
const targetsPerContainer = 10

func StartScan(c *fiber.Ctx) error {
	log.Println("Starting scan...")

	// First find total waiting scan count
	var zapWaitingScanCount = model.GetWaitingZAPTargetDomainCount()
	log.Println("ZAP Waiting Scan Count:", zapWaitingScanCount)
	var openvasWaitingScanCount = model.GetWaitingOpenVASTargetDomainCount()
	log.Println("OpenVAS Waiting Scan Count:", openvasWaitingScanCount)
	var nucleiWaitingScanCount = model.GetWaitingNucleiTargetDomainCount()
	log.Println("Nuclei Waiting Scan Count:", nucleiWaitingScanCount)

	dockerManager := hunterbounter_docker.NewDockerManager()

	// Function to remove a container by ID
	removeContainer := func(containerID string) {
		err := dockerManager.RemoveContainer(containerID)
		if err != nil {
			log.Println("Error while removing container", containerID, ":", err)
		} else {
			log.Println("Successfully removed container", containerID)
		}
	}

	// Function to adjust containers for an agent
	adjustContainers := func(agentName string, waitingCount int) {
		log.Printf("Adjusting containers for agent %s...", agentName)
		runningContainers, err := dockerManager.ListRunningContainersViaImageName(agentName)
		if err != nil {
			log.Printf("Error while listing running containers for %s: %v", agentName, err)
			return
		}
		runningCount := len(runningContainers)
		requiredContainers := (waitingCount + targetsPerContainer - 1) / targetsPerContainer // Round up to the nearest whole number
		log.Printf("%s - Running: %d, Required: %d", agentName, runningCount, requiredContainers)

		if requiredContainers > maxContainers {
			requiredContainers = maxContainers
		}

		if runningCount > requiredContainers {
			containersToStop := runningCount - requiredContainers
			log.Printf("Stopping %d containers for %s", containersToStop, agentName)
			for i := 0; i < containersToStop; i++ {
				removeContainer(runningContainers[i])
			}
		} else if runningCount < requiredContainers {
			containersToStart := requiredContainers - runningCount
			log.Printf("Starting %d containers for %s", containersToStart, agentName)
			for i := 0; i < containersToStart; i++ {
				var imageName string
				switch agentName {
				case DOCKER_DOMAIN_IMAGE_NAME_ZAP:
					imageName = "hunter_bounter_zapv1"
				case DOCKER_DOMAIN_IMAGE_NAME_OPENVAS:
					imageName = "hunterbounter:openvasv1"
				case DOCKER_DOMAIN_IMAGE_NAME_NUCLEI:
					imageName = "hunterbounter:nucleiv1"
				}

				if agentName == DOCKER_DOMAIN_IMAGE_NAME_ZAP {
					user := "root"
					port := "5002:5002"
					dns := "1.1.1.1"
					log.Println("Starting ZAP container...")
					containerID, err := dockerManager.RunContainer(imageName, user, port, dns)
					log.Println("Running ZAP container ID:", containerID)
					if err != nil {
						log.Println("Error while starting ZAP container:", err)
						return
					}
				} else {
					log.Printf("Starting container for %s...", agentName)
					containerID, err := dockerManager.RunContainerNoCommand(imageName)
					log.Printf("Running %s container ID: %s", agentName, containerID)
					if err != nil {
						log.Printf("Error while starting %s container: %v", agentName, err)
						return
					}
				}
			}
		}
	}

	// Adjust containers for each agent
	adjustContainers(DOCKER_DOMAIN_IMAGE_NAME_ZAP, zapWaitingScanCount)
	adjustContainers(DOCKER_DOMAIN_IMAGE_NAME_OPENVAS, openvasWaitingScanCount)
	adjustContainers(DOCKER_DOMAIN_IMAGE_NAME_NUCLEI, nucleiWaitingScanCount)

	// Log the final state
	var zapLen, _ = dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME_ZAP)
	var openvasLen, _ = dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME_OPENVAS)
	var nucleiLen, _ = dockerManager.ListRunningContainersViaImageName(DOCKER_DOMAIN_IMAGE_NAME_NUCLEI)

	log.Println("Final state after adjustments:")
	log.Println("ZAP running containers:", zapLen)
	log.Println("OpenVAS running containers:", openvasLen)
	log.Println("Nuclei running containers:", nucleiLen)

	// Type 1 Domain
	// Status 0 Added
	log.Println("Updating target statuses...")
	model.UpdateTargetsStatus(0, 1, 1)
	model.UpdateTargetsStatus(0, 2, 1)
	model.UpdateTargetsStatus(0, 3, 1)

	log.Println("Scan started successfully.")
	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan started", nil))
}

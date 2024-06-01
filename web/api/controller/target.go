package controller

import (
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"log"
	"reflect"
	"sync"
)

const MAX_RUNNING_SCAN_COUNT = 10

// enum docker type
type DockerType int

const (
	ZAP = iota + 1
	OpenVAS
)

type TargetRequest struct {
	TotalRunningScanCount int `json:"total_running_scan_count"`
	DockerType            int `json:"docker_type"`
}

type TargetResponse struct {
	Targets []string `json:"targets"`
}

var targetMutex = &sync.Mutex{}

func GetTargets(c *fiber.Ctx) error {

	log.Println("Get Targets")

	// We Use Mutex because we don't want to get the same target multiple times
	targetMutex.Lock()
	defer targetMutex.Unlock()

	var request TargetRequest

	if err := c.BodyParser(&request); err != nil {
		log.Println(hunterbounter_json.ToString(request))
		log.Println("Error parsing request", err)

		var testMap map[string]interface{}
		err := c.BodyParser(&testMap)
		if err != nil {
			log.Println("Error parsing request", err)
		}
		log.Println("test map -> ", hunterbounter_json.ToString(testMap))
		log.Println(reflect.TypeOf(testMap["total_running_scan_count"]))
		log.Println(reflect.TypeOf(testMap["docker_type"]))
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error parsing request", nil))
	}

	targetsResponse, err := database.Select("targets", map[string]interface{}{
		"status": "1", // waiting
		"type":   "1", // 1 domain, 2 ipv4
	})

	if err != nil {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error getting targets", nil))
	}

	var targets []string
	var maxCount int
	for _, target := range targetsResponse {

		targets = append(targets, target["value"].(string))
		maxCount++
		if maxCount >= MAX_RUNNING_SCAN_COUNT {
			break
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Targets", TargetResponse{Targets: targets}))

}

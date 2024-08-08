package controller

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"html"
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/utils"
	"log"
)

func CheckNucleiRecordIsExist(elemRecord map[string]interface{}) bool {
	var whereCondition = map[string]interface{}{"machine_id": elemRecord["machine_id"], "full_json": elemRecord["full_json"], "task_id": elemRecord["task_id"], "timestamp": elemRecord["timestamp"]}

	// Check if the record is exist
	dbRecords, err := database.Select("nuclei_scan_results", whereCondition)

	if err != nil {
		log.Println("Error while checking record is exist", err)
		return false
	}

	if len(dbRecords) > 0 {
		return true
	}
	return false

}

func ScanResultNucleiPOST(c *fiber.Ctx) error {

	log.Println("Scan Result Post (Nuclei)")

	var scanResult []map[string]interface{}
	err := c.BodyParser(&scanResult)
	if err != nil {
		log.Println("Request Body:", string(c.Body()))

		if string(c.Body()) == "" {
			log.Println("Empty Body")
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Empty Body", nil))
		}

		log.Println("Error parsing request", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error parsing request", nil))
	}

	for index, elemRecord := range scanResult {
		for key, value := range elemRecord {
			if key == "info" {
				continue
			}
			elemRecord[key] = html.EscapeString(utils.SafeEscapeString(value))
		}

		jsonData, err := json.Marshal(elemRecord)
		if err != nil {
			log.Fatal(err)
		}

		encodedData := base64.StdEncoding.EncodeToString(jsonData)
		elemRecord["full_json"] = encodedData
		log.Println("record index", index)
		if CheckNucleiRecordIsExist(elemRecord) {
			log.Println("Record is exist")
			continue
		}

		var info = elemRecord["info"].(interface{})
		log.Println(hunterbounter_json.ToStringBeautify(info))

		log.Println("Record is not exist", hunterbounter_json.ToString(elemRecord))
		var inserData = map[string]interface{}{
			"ip":         elemRecord["ip"],
			"name":       utils.SafeEscapeString(elemRecord["info"].(interface{}).(map[string]interface{})["name"]),
			"severity":   utils.SafeEscapeString(elemRecord["info"].(interface{}).(map[string]interface{})["severity"]),
			"host":       elemRecord["host"],
			"machine_id": elemRecord["machine_id"],
			"task_id":    elemRecord["task_id"],
			"timestamp":  elemRecord["timestamp"],
			"full_json":  encodedData,
		}

		// Insert the record
		_, err = database.Insert("nuclei_scan_results", inserData, false)
		if err != nil {
			log.Println("insert_data -> ", hunterbounter_json.ToString(inserData))
			log.Println("elemRecord -> ", hunterbounter_json.ToString(elemRecord))
			log.Println("Error while inserting record", err)
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while inserting record", nil))
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan Result Received", nil))
}

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

func CheckRecordIsExist(elemRecord map[string]interface{}) bool {
	var whereCondition = map[string]interface{}{"machine_id": elemRecord["machine_id"], "url": elemRecord["url"], "zap_id": elemRecord["id"]}

	// Check if the record is exist
	dbRecords, err := database.Select("zap_scan_results", whereCondition)

	if err != nil {
		log.Println("Error while checking record is exist", err)
		return false
	}

	if len(dbRecords) > 0 {
		return true
	}
	return false

}

func ScanResultPOST(c *fiber.Ctx) error {

	log.Println("Scan Result Post (OpenVAS)")

	var scanResult []map[string]interface{}
	err := c.BodyParser(&scanResult)
	if err != nil {
		log.Println("Request Body:", string(c.Body()))

		log.Println("Error parsing request", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error parsing request", nil))
	}

	// check if the record is exist
	for _, elemRecord := range scanResult {
		if CheckRecordIsExist(elemRecord) {
			//log.Println("Record is exist")
			continue
		}

		//elemRecord["description"] = html.EscapeString(utils.SafeEscapeString(elemRecord["description"])) // deprecated
		//elemRecord["solution"] = html.EscapeString(utils.SafeEscapeString(elemRecord["solution"])) // deprecated
		//elemRecord["other"] = html.EscapeString(utils.SafeEscapeString(elemRecord["other"])) // deprecated
		//json_string := hunterbounter_json.ToString(elemRecord)

		for key, value := range elemRecord {
			elemRecord[key] = html.EscapeString(utils.SafeEscapeString(value))
		}

		jsonData, err := json.Marshal(elemRecord)
		if err != nil {
			log.Fatal(err)
		}
		encodedData := base64.StdEncoding.EncodeToString(jsonData)

		//log.Println("Record is not exist", hunterbounter_json.ToString(elemRecord))
		var inserData = map[string]interface{}{
			"url":         elemRecord["url"],
			"zap_id":      elemRecord["id"],
			"risk":        elemRecord["risk"],
			"description": elemRecord["description"],
			"solution":    elemRecord["solution"],
			"other_info":  elemRecord["other"].(string),
			"reference":   elemRecord["reference"],
			"cwe_id":      elemRecord["cweid"],
			"wasc_id":     elemRecord["wascid"],
			"machine_id":  elemRecord["machine_id"],
			"full_json":   encodedData,
		}

		// Insert the record
		_, err = database.Insert("zap_scan_results", inserData, false)
		if err != nil {
			log.Println("insert_data -> ", hunterbounter_json.ToString(inserData))
			log.Println("elemRecord -> ", hunterbounter_json.ToString(elemRecord))
			log.Println("Error while inserting record", err)
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while inserting record", nil))
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan Result Received", nil))
}

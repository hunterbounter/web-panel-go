package controller

import (
	"github.com/gofiber/fiber/v2"
	"html"
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
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

	log.Println("Scan Result Post")

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
		//log.Println("Record is not exist", hunterbounter_json.ToString(elemRecord))
		var inserData = map[string]interface{}{
			"url":         elemRecord["url"],
			"zap_id":      elemRecord["id"],
			"risk":        elemRecord["risk"],
			"description": html.EscapeString(elemRecord["description"].(string)),
			"solution":    html.EscapeString(elemRecord["solution"].(string)),
			"other_info":  elemRecord["other"].(string),
			"reference":   elemRecord["reference"],
			"cwe_id":      elemRecord["cweid"],
			"wasc_id":     elemRecord["wascid"],
			"machine_id":  elemRecord["machine_id"],
		}

		// Insert the record
		_, err := database.Insert("zap_scan_results", inserData, false)
		if err != nil {
			log.Println("Error while inserting record", err)
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while inserting record", nil))
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan Result Received", nil))
}

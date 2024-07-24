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

func CheckOpenVASRecordIsExist(elemRecord map[string]interface{}) bool {
	var whereCondition = map[string]interface{}{"machine_id": elemRecord["machine_id"], "result_id": elemRecord["Result ID"]}

	// Check if the record is exist
	dbRecords, err := database.Select("openvas_scan_result", whereCondition)

	if err != nil {
		log.Println("Error while checking record is exist", err)
		return false
	}

	if len(dbRecords) > 0 {
		return true
	}
	return false

}

func ScanResultOpenVASPOST(c *fiber.Ctx) error {

	log.Println("Scan Result Post (OpenVAS)")

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

	log.Println("Scan Result Received")
	for index, elemRecord := range scanResult {
		log.Println("record index", index)
		if CheckOpenVASRecordIsExist(elemRecord) {
			log.Println("Record is exist")
			continue
		}

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
			"ip":                             elemRecord["IP"],
			"hostname":                       elemRecord["Hostname"],
			"port":                           elemRecord["Port"],
			"port_protocol":                  elemRecord["Port Protocol"],
			"cvss":                           elemRecord["CVSS"],
			"severity":                       elemRecord["Severity"],
			"qod":                            elemRecord["QoD"],
			"solution":                       elemRecord["Solution"],
			"solution_type":                  elemRecord["Solution Type"],
			"nvt_name":                       elemRecord["NVT Name"],
			"summary":                        elemRecord["Summary"],
			"nvt_oid":                        elemRecord["NVT OID"],
			"specific_result":                elemRecord["Specific Result"],
			"task_id":                        elemRecord["Task ID"],
			"task_name":                      elemRecord["Task Name"],
			"timestamp":                      elemRecord["Timestamp"],
			"result_id":                      elemRecord["Result ID"],
			"impact":                         elemRecord["Impact"],
			"affected_software_or_os":        elemRecord["Affected Software/OS"],
			"vulnerability_insight":          elemRecord["Vulnerability Insight"],
			"vulnerability_detection_method": elemRecord["Vulnerability Detection Method"],
			"product_detection_result":       elemRecord["Product Detection Result"],
			"machine_id":                     elemRecord["machine_id"],
			//"bids":                           elemRecord["Summary"],
			//"certs":                          elemRecord["Summary"],
			//"other_reference":                elemRecord["Summary"],
			"full_json": encodedData,
		}

		// Insert the record
		_, err = database.Insert("openvas_scan_result", inserData, false)
		if err != nil {
			log.Println("insert_data -> ", hunterbounter_json.ToString(inserData))
			log.Println("elemRecord -> ", hunterbounter_json.ToString(elemRecord))
			log.Println("Error while inserting record", err)
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while inserting record", nil))
		}

	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan Result Received", nil))
}

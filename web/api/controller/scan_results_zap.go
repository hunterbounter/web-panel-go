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
	"path/filepath"
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

func ScanResultSavePDF(c *fiber.Ctx) error {
	log.Println("Scan Result Save PDF")

	uuidFileName := c.FormValue("uuid_file_name")
	log.Println("UUID File Name: ", uuidFileName)
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error while getting file", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while getting file", nil))
	}
	var destination string

	if file != nil {
		fileExtension := filepath.Ext(file.Filename)
		// md5
		var fileMD5, err = utils.GenerateMD5(file)
		if err != nil {
			log.Println("Error while generating MD5", err)
			return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Error while generating MD5", nil))
		}
		log.Println("File MD5: ", fileMD5)

		uniqueFileName := utils.GenerateUUID() + fileExtension

		destination = "/uploads/" + uniqueFileName
		err = c.SaveFile(file, utils.RunningDir()+destination)
		if err != nil {
			log.Println("Error: ", err)
			return c.Status(fiber.StatusOK).JSON(hunterbounter_response.HunterBounterResponse(false, "File upload failed", nil))
		}

		var updateData = map[string]interface{}{
			"pdf_file": destination,
			"status":   3, // completed
		}
		var whereCondition = map[string]interface{}{"uuid_file_name": uuidFileName}

		update, err := database.Update("mobile_application_results", updateData, whereCondition)
		if err != nil {
			log.Println("Error while updating record", err)
			return err
		}
		log.Println("Update Result: ", update)
	} else {
		log.Println("File is nil")
	}
	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Scan Result Received", nil))
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
			"name":        elemRecord["name"],
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

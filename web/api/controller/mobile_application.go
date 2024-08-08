package controller

import (
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/utils"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
	"path/filepath"
)

func MobileScanGET(c *fiber.Ctx) error {

	var mobileScanList = model.GetMobileScanList()

	return RenderTemplate(c, "panel/mobile_application", fiber.Map{
		"Title":          "Mobile Application Scan",
		"MobileScanList": mobileScanList,
	})
}
func UploadMobileAppFile(c *fiber.Ctx) error {
	var scanName string
	scanName = c.FormValue("scan_name")

	if scanName == "" {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Scan name is required", nil))

	}

	file, err := c.FormFile("file")

	if err != nil {
		log.Println("File Upload Error: ", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "File upload failed", nil))
	}

	var destination string

	uniqueFileName := utils.GenerateUUID()
	if file != nil {
		// Check if the file is an APK file
		fileExtension := filepath.Ext(file.Filename)

		uniqueFileName = uniqueFileName + fileExtension

		destination = "/uploads/" + uniqueFileName
		err = c.SaveFile(file, utils.RunningDir()+destination)
		if err != nil {
			log.Println("Error: ", err)
			return c.Status(fiber.StatusOK).JSON(hunterbounter_response.HunterBounterResponse(false, "File upload failed", nil))
		}
	}

	var insertData = map[string]interface{}{
		"scan_name":      scanName,
		"file_name":      file.Filename,
		"status":         1,
		"app_dir":        destination,
		"uuid_file_name": uniqueFileName,
	}

	// Save to database
	_, err = database.Insert("mobile_application_results", insertData, false)
	if err != nil {
		log.Println("Error: ", err)
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "File upload failed", nil))
	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "File uploaded successfully", nil))

}

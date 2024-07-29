package model

import (
	"hunterbounter.com/web-panel/pkg/database"
	"log"
)

func GetTotalReportCount() int {

	dbRecords, err := database.Select("zap_scan_results", map[string]interface{}{})
	if err != nil {
		log.Println("Error while checking record is exist", err)
	}

	dbRecords2, err2 := database.Select("openvas_scan_result", map[string]interface{}{})
	if err2 != nil {
		log.Println("Error while checking record is exist", err2)
	}

	dbRecords3, err3 := database.Select("nuclei_scan_results", map[string]interface{}{})
	if err2 != nil {
		log.Println("Error while checking record is exist", err3)
	}

	return len(dbRecords) + len(dbRecords2) + len(dbRecords3)
}

package model

import "hunterbounter.com/web-panel/pkg/database"

func GetTotalReportCount() int {

	dbRecords, err := database.Select("zap_scan_results", map[string]interface{}{})
	if err != nil {
		return 0
	}
	return len(dbRecords)
}

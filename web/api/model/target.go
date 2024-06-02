package model

import "hunterbounter.com/web-panel/pkg/database"

func SaveTarget(target string, elemType int, status int) {

	database.Insert("targets", map[string]interface{}{"value": target, "type": elemType, "status": status}, false)

}

func GetWaitingTargetDomainCount() int {
	var whereCondition = map[string]interface{}{"status": 0, "type": 1} // 0 Added 1 Waiting 2 Completed --- Type 1 Domain 2 IP
	dbRecords, err := database.Select("targets", whereCondition)
	if err != nil {
		return 0
	}
	return len(dbRecords)
}

func GetTotalScannedTargetCount() int {
	var whereCondition = map[string]interface{}{"status": 2} // 0 Added 1 Waiting 2 Completed --- Type 1 Domain 2 IP
	dbRecords, err := database.Select("targets", whereCondition)
	if err != nil {
		return 0
	}
	return len(dbRecords)
}

func UpdateTargetsStatus(whereStatus int, whereType int, updateStatus int) {
	updateData := map[string]interface{}{"status": updateStatus}
	whereData := map[string]interface{}{"status": whereStatus, "type": whereType}
	database.Update("targets", updateData, whereData)

	// String : update targets set status = 1 where status = 0 and type = 1
}

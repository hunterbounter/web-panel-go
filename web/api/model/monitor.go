package model

import "hunterbounter.com/web-panel/pkg/database"

func GetMonitorData() []map[string]interface{} {
	var viewMonitorSql = `
	select * from VIEW_MONITOR order by hostname asc
	 
`

	dbRecords, err := database.ExecuteSql(viewMonitorSql)
	if err != nil {
		return nil
	}
	return dbRecords

}

func RemoveMonitorData(hostname string) bool {
	var deleteMonitorSql = `
	delete from machine_monitor where hostname = '` + hostname + `'
`

	_, err := database.ExecuteSql(deleteMonitorSql)
	if err != nil {
		return false
	}
	return true

}

package model

import "hunterbounter.com/web-panel/pkg/database"

func GetMobileScanList() []map[string]interface{} {

	var viewMonitorSql = `
	select
	id,
	scan_name,
	file_name,
	package_name,
	version_name,
	md5,
	app_dir,
	pdf_file,
	case when status = 1 then 'Waiting' when status = 2 then 'Scanning' when status = 3 then 'Completed' end as status,
	create_date
from
	mobile_application_results
	 
`

	dbRecords, err := database.ExecuteSql(viewMonitorSql)
	if err != nil {
		return nil
	}
	return dbRecords

}

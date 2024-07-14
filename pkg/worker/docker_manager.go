package worker

import (
	"hunterbounter.com/web-panel/pkg/date_util"
	"hunterbounter.com/web-panel/pkg/docker"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	"hunterbounter.com/web-panel/web/api/model"
	"log"
)

const REMOVE_FROM_LIST_TIMEOUT = 300 // 5 minutes in seconds

/*
	We are creating a map, if the scan_count is 0, it will be deleted after 5 minutes.
*/

var workerMap = map[string]string{}

func workerDockerManager() {

	/*
		First of all, if the Last Seen time is more than 5 minutes, kill the docker processes and delete from the monitor table.
	*/

	var monitorData = model.GetMonitorData()

	for _, monitor := range monitorData {
		log.Println("workerMap : ", hunterbounter_json.ToStringBeautify(workerMap))
		var hostname = monitor["hostname"].(string)
		var lastSeen = monitor["last_seen"].(string) // 14.07.2024 19:51:05
		var activeScanCount = monitor["active_scan_count"].(int64)
		var currentTime = date_util.DateYYYYMMDDHH24MISS()

		// if active scan count is 0
		if activeScanCount == 0 {
			if _, ok := workerMap[hostname]; !ok {
				log.Println("Adding to list for ", hostname, " Last Seen : ", lastSeen, " Current Time : ", currentTime)
				workerMap[hostname] = lastSeen
			}

			// mapdeki last seen süresi ile şuanki zaman arasındaki farkı alıyoruz.
			var mapLastSeen = workerMap[hostname]

			if mapLastSeen == "" {
				continue
			}
			diff, _ := date_util.DateDiff(mapLastSeen, currentTime)
			log.Println("Diff : ", diff, " for ", hostname)
			if diff > REMOVE_FROM_LIST_TIMEOUT {
				log.Println("Removing from list for ", hostname, " Last Seen : ", lastSeen, " Current Time : ", currentTime, " Diff : ", diff)
				model.RemoveMonitorData(hostname)
				err := docker.NewDockerManager().RemoveContainer(hostname)
				if err != nil {
					log.Println("Removing from list failed for ", hostname, " error : ", err)
				}
				delete(workerMap, hostname)
			}

			continue
		} else {
			delete(workerMap, hostname)
		}

	}

}

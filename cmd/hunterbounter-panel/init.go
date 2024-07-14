package main

import (
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	hunterbounter_session "hunterbounter.com/web-panel/pkg/session"
	"hunterbounter.com/web-panel/pkg/utils"
	"hunterbounter.com/web-panel/pkg/worker"
	"log"
)

func InitSys() {
	// ENV
	var sysEnv = InitEnv()

	// init worker
	worker.Init()
	// init web session
	hunterbounter_session.InitSession()

	// init database

	log.Println("sys env ", hunterbounter_json.ToStringBeautify(sysEnv))

	postgres_info := database.PostgreSQLInfo{
		Username:  sysEnv.DBUsername,
		Password:  sysEnv.DBPassword,
		Hostname:  sysEnv.DBHost,
		Port:      utils.AnyToInt(sysEnv.DBPort),
		Name:      sysEnv.DBName,
		Parameter: sysEnv.DBParams,
	}

	info := database.Info{
		Type:       database.Type(sysEnv.DBType),
		PostgreSQL: postgres_info,
	}
	database.Connect(info)

}

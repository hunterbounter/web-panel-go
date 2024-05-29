package main

import (
	"fmt"
	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/env"
	"hunterbounter.com/web-panel/web/api/bootstrap"
	"log"
)

func init() {
	postgres_info := database.PostgreSQLInfo{
		Username:  "zap_user",
		Password:  "your_password",
		Hostname:  "127.0.0.1",
		Port:      5432,
		Name:      "postgres",
		Parameter: "sslmode=disable",
	}

	info := database.Info{
		Type:       "PostgreSQL",
		PostgreSQL: postgres_info,
	}
	database.Connect(info)
}

func main() {
	log.Println("Running HunterBounter Web Panel...")
	panelApp := bootstrap.HunterBounterWeb()
	err := panelApp.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "0.0.0.0"), env.GetEnv("APP_PORT", "9000")))
	if err != nil {
		log.Println("Failed to start panel web : ", err)
		return
	}
}

package main

import (
	"github.com/joho/godotenv"
	"hunterbounter.com/web-panel/pkg/utils"
	"log"
	"os"
)

// ENVIRONMENT VARIABLES
type EnvirontmentVariables struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBParams   string
	DBType     string
}

func InitEnv() *EnvirontmentVariables {
	err := godotenv.Load(utils.RunningDir() + ".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &EnvirontmentVariables{
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBParams:   os.Getenv("DB_PARAMS"),
		DBType:     os.Getenv("DB_TYPE"),
	}
}

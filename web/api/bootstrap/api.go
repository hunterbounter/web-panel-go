package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"hunterbounter.com/web-panel/pkg/router"
	"log"
	"os"
	"path/filepath"
)

func RunningDir() string {
	/*if os.Getenv("USER") == "eminsargin" {

		return "/Users/eminsargin/Desktop/PROJECTS/github.com/github.com/esargin/hunterbounter.com/web-panel"
	}*/
	path, err := os.Executable()
	if err != nil {
		return ""
	}

	// the executable directory
	exPath := filepath.Dir(path)

	exPath = exPath + "/../../"

	return exPath
}

func HunterBounterWeb() *fiber.App {
	//env.SetupEnvFile()

	app := fiber.New(fiber.Config{BodyLimit: 500 * 1024 * 1024})

	app.Use(recover.New())

	app.Use(helmet.New(
		helmet.Config{
			XSSProtection: "1; mode=block",
		}))

	// Set the log output
	log.SetOutput(os.Stdout)

	router.InstallRouter(app)

	return app
}

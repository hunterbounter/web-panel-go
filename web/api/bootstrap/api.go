package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"hunterbounter.com/web-panel/pkg/router"

	"log"
	"os"
	"path/filepath"
)

func RunningDir() string {
	if os.Getenv("USER") == "eminsargin" {

		return "/Users/eminsargin/Desktop/PROJECTS/github.com/hunterbounter.com/web-panel-go"
	}
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
	engine := html.New(RunningDir()+"/views", ".html")

	app := fiber.New(fiber.Config{Views: engine, BodyLimit: 500 * 1024 * 1024})

	app.Use(recover.New())

	app.Use(helmet.New(
		helmet.Config{
			XSSProtection: "1; mode=block",
		}))

	// Set the log output
	log.SetOutput(os.Stdout)

	log.Println(RunningDir() + "/views/statics")
	app.Static("/", RunningDir()+"/views/statics")

	app.Get("/sys", monitor.New())

	router.InstallRouter(app)

	return app
}

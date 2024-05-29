package router

import "github.com/gofiber/fiber/v2"

func InstallRouter(app *fiber.App) {
	setup(app, NewHunterBounterApiRouter())
}
func setup(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}

package controller

import (
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/web/api/model"
)

func GetTargetsView(c *fiber.Ctx) error {
	var targetList = model.GetTargetList()
	return c.Render("panel/targets", fiber.Map{
		"Title":   "Targets",
		"Targets": targetList,
	})
}

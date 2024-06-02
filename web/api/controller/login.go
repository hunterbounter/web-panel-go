package controller

import (
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	hunterbounter_session "hunterbounter.com/web-panel/pkg/session"
)

func LoginGet(c *fiber.Ctx) error {
	return RenderTemplate(c, "panel/login", fiber.Map{})
}

type LoginRequest struct {
	MailAddress string `json:"mail_address"`
	Password    string `json:"password"`
}

func validateLogin(loginRequest LoginRequest) *hunterbounter_response.HunterResponse {
	if loginRequest.MailAddress == "" {
		return hunterbounter_response.HunterBounterResponse(false, "Mail adresi boş olamaz.", nil)
	}
	if loginRequest.Password == "" {
		return hunterbounter_response.HunterBounterResponse(false, "Şifre boş olamaz.", nil)
	}
	return nil

}

func LoginPost(c *fiber.Ctx) error {
	var mailAddress = c.FormValue("mail_address")
	var password = c.FormValue("password")

	var loginRequest LoginRequest

	loginRequest.MailAddress = mailAddress
	loginRequest.Password = password

	validationResponse := validateLogin(loginRequest)

	if validationResponse != nil {
		return c.Status(fiber.StatusOK).JSON(validationResponse)
	}

	if loginRequest.MailAddress == "admin" && loginRequest.Password == "admin" {
		hunterbounter_session.SetSessionValue(c, "userId", "admin")
	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Login successful", nil))

}

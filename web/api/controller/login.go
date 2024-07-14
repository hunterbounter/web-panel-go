package controller

import (
	"github.com/gofiber/fiber/v2"
	"log"

	"hunterbounter.com/web-panel/pkg/database"
	"hunterbounter.com/web-panel/pkg/date_util"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/passhash"
	hunterbounter_session "hunterbounter.com/web-panel/pkg/session"
)

func LoginGet(c *fiber.Ctx) error {
	log.Println("LoginGet")
	return RenderTemplate(c, "panel/login", fiber.Map{})
}

type LoginRequest struct {
	MailAddress string `json:"mail_address"`
	Password    string `json:"password"`
}

func validateLogin(loginRequest LoginRequest) *hunterbounter_response.HunterResponse {
	if loginRequest.MailAddress == "" {
		return hunterbounter_response.HunterBounterResponse(false, "Mail Address Cannot Be Empty", nil)
	}
	if loginRequest.Password == "" {
		return hunterbounter_response.HunterBounterResponse(false, "Password Cannot Be Empty", nil)
	}
	return nil

}

// Database validation
func validateUserInformationViaDB(loginRequest LoginRequest, dbResponse []map[string]interface{}) *hunterbounter_response.HunterResponse {

	if len(dbResponse) == 0 {
		return hunterbounter_response.HunterBounterResponse(false, "Username or Password is incorrect", nil)
	}
	var dbHashedPassword = dbResponse[0]["password"].(string)

	// Check if the password is correct
	if !passhash.MatchString(dbHashedPassword, loginRequest.Password) {
		return hunterbounter_response.HunterBounterResponse(false, "Username or Password is incorrect", nil)
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

	var whereData = map[string]interface{}{
		"mail": loginRequest.MailAddress,
	}
	dbResponse, err := database.Select("users", whereData)
	log.Println("dbResponse", dbResponse)

	if err != nil {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "Undefined Error", nil))
	}

	// Set session
	var sessionData = hunterbounter_session.SessionData{
		UserID:       dbResponse[0]["id"].(string),
		Username:     dbResponse[0]["username"].(string),
		IpAddress:    GetRealIp(c),
		LoginTime:    date_util.DateYYYYMMDDHH24MISS(),
		LastActivity: date_util.DateYYYYMMDDHH24MISS(),
	}

	hunterbounter_session.SetSessionValue(c, hunterbounter_session.SESSION_KEY, &sessionData)

	// Database validation
	dbValidationResponse := validateUserInformationViaDB(loginRequest, dbResponse)

	if dbValidationResponse != nil {
		return c.Status(fiber.StatusOK).JSON(dbValidationResponse)
	}

	var updateData = map[string]interface{}{
		"last_login_date": database.CurrentTime(),
		"last_login_ip":   GetRealIp(c),
	}

	// Update last login date and ip
	_, err = database.Update("users", updateData, whereData)
	if err != nil {
		log.Println("Update Error: ", err)
		return err
	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Login successful", nil))

}

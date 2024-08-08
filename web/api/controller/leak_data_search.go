package controller

import (
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"hunterbounter.com/web-panel/pkg/hunterbounter_response"
	"hunterbounter.com/web-panel/pkg/utils"
	"log"
	"time"
)

const WHITE_INTEl_CUSTOMER_LEAK_ENDPOINT = "https://whiteintel.io/api/public/get_customer_leaks_public.php"
const WHITE_INTEl_API_KEY = "YOUR_API_KEY"

type WhiteIntelRequest struct {
	Apikey            string `json:"apikey"`
	Query             string `json:"query"`
	IncludeSystemInfo int    `json:"include_system_info"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
}

type LeakDataResponse struct {
	TotalLeaks int `json:"total_leaks"`
	Data       []struct {
		Url      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"data"`
}
type LeakDataRequest struct {
	Query string `json:"query"`
}

func LeakDataSearchPost(c *fiber.Ctx) error {
	client := resty.New()

	var leakdataRequest LeakDataRequest

	if err := c.BodyParser(&leakdataRequest); err != nil {
		log.Println("Error parsing body: ", err)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid Request",
		})

	}
	log.Println("Request: ", leakdataRequest)
	if !utils.IsValidDomainNormal(leakdataRequest.Query) {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid Domain",
		})
	}

	var whiteIntelRequest WhiteIntelRequest
	whiteIntelRequest.Apikey = WHITE_INTEl_API_KEY
	whiteIntelRequest.Query = leakdataRequest.Query
	whiteIntelRequest.IncludeSystemInfo = 0
	whiteIntelRequest.StartDate = "2020-01-01"
	whiteIntelRequest.EndDate = "2025-01-01"

	var leakDataResponse LeakDataResponse

	client.SetTimeout(60 * time.Second)
	client.Debug = true

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "HunterBounter Version/1.0").
		SetBody(whiteIntelRequest).
		SetResult(&leakDataResponse).
		//SetDebug(true).
		Post(WHITE_INTEl_CUSTOMER_LEAK_ENDPOINT)

	if err != nil {
		log.Println("Error making request: ", err)
	}
	log.Println("Response: ", leakDataResponse)

	for index, _ := range leakDataResponse.Data {
		leakDataResponse.Data[index].Password = "***********"
		leakDataResponse.Data[index].Username = "***********"

	}
	log.Println("Response: ", leakDataResponse)
	if len(leakDataResponse.Data) > 25 {
		leakDataResponse.Data = leakDataResponse.Data[0:13]
	}
	if len(leakDataResponse.Data) == 0 {
		return c.JSON(hunterbounter_response.HunterBounterResponse(false, "No Data Found", nil))
	}

	return c.JSON(hunterbounter_response.HunterBounterResponse(true, "Success", leakDataResponse))
}
func LeakDataSearchGET(c *fiber.Ctx) error {

	return RenderTemplate(c, "panel/leak_data_search", fiber.Map{
		"Title": "Leak Data Search",
	})

}

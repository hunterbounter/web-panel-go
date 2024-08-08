package controller

import (
	"github.com/go-resty/resty/v2"
	"hunterbounter.com/web-panel/pkg/hunterbounter_json"
	"log"
	"testing"
	"time"
)

func TestLeakData(t *testing.T) {
	client := resty.New()

	var leakDataRequest LeakDataRequest

	leakDataRequest.Query = "defcon.org"
	var whiteIntelRequest WhiteIntelRequest
	whiteIntelRequest.Apikey = WHITE_INTEl_API_KEY
	whiteIntelRequest.Query = leakDataRequest.Query
	whiteIntelRequest.IncludeSystemInfo = 0
	whiteIntelRequest.StartDate = "2020-01-01"
	whiteIntelRequest.EndDate = "2025-01-01"

	var leakDataResponse LeakDataResponse

	client.SetTimeout(60 * time.Second)
	client.Debug = true
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
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
	log.Println("Response: ", hunterbounter_json.ToStringBeautify(leakDataResponse))
}

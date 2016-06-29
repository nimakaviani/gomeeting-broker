package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nimakaviani/gomeeting-broker/models"
	"github.com/nimakaviani/gomeeting-broker/utils"
	"github.com/pivotal-golang/lager"
)

type handler struct {
	gCalendar utils.GCalendar
}

func NewHandler(gCalendar utils.GCalendar) handler {
	return handler{
		gCalendar: gCalendar,
	}
}

func (h handler) Alexa(writer http.ResponseWriter, request *http.Request) {
	logger := lager.NewLogger("alexa")
	alexaRequest := models.AlexaRequest{}
	err := json.NewDecoder(request.Body).Decode(&alexaRequest)

	alexaResp := models.NewAlexaResponse()

	fmt.Printf("request %#v\n", alexaRequest)
	startTime, duration, err := utils.Parse(alexaRequest)

	h.gCalendar.FindRoom(startTime, duration)

	if err != nil {
		logger.Error("failed-parse-duration", err)
		alexaResp.OutputSpeech("I could not understand your request")
	} else {
		alexaResp.OutputSpeech(utils.PrepareResponse(startTime, duration))
		if err != nil {
			logger.Error("failed-prepare-response", err)
		}
	}

	json, err := alexaResp.String()
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.Write(json)

	logger.Info("request: " + fmt.Sprintf("%#v", alexaRequest))
}

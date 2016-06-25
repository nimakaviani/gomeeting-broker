package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nimakaviani/gomeeting-broker/models"
	"github.com/pivotal-golang/lager"
)

type handler struct {
}

func NewHandler() handler {
	return handler{}
}

func (h handler) Alexa(writer http.ResponseWriter, request *http.Request) {
	logger := lager.NewLogger("alexa")
	alexaRequest := models.AlexaRequest{}
	err := json.NewDecoder(request.Body).Decode(&alexaRequest)

	startTime, duration, err := parseDuration(alexaRequest)
	if err != nil {
		logger.Error("failed-parse-duration", err)
	}

	alexaResp := models.NewAlexaResponse()
	alexaResp.OutputSpeech(prepareResponse(startTime, duration))
	json, err := alexaResp.String()
	if err != nil {
		logger.Error("failed-prepare-response", err)
	}

	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.Write(json)

	logger.Info("request: " + fmt.Sprintf("%#v", alexaRequest))
}

func prepareResponse(startTime time.Time, duration time.Duration) string {
	return "Hello world!"
}

func parseDuration(alexaRequest models.AlexaRequest) (time.Time, time.Duration, error) {
	return time.Now(), 1 * time.Hour, nil
}

func parseTime(alexaRequest models.AlexaRequest) time.Duration {
	return 10 * time.Hour
}

func humanizeLength(length int) (int, string) {
	switch {
	case length <= 60:
		return 1, "minute"
	case length < 3600:
		return length / 60, "minutes"
	case length > 3600:
		return length / 3600, "hours"
	default:
		return 1, "hour"
	}
}

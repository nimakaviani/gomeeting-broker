package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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

	startTime, duration, err := parse(alexaRequest)
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
	return fmt.Sprintf("Hello world! %#v", startTime)
}

func parse(alexaRequest models.AlexaRequest) (time.Time, time.Duration, error) {
	// startInterval := parseTime(alexaRequest.Request.Intent.Slots["StartTime"])
	// duration := parseDuration(alexaRequest.Request.Intent.Slots["Length"])
	// println("startInterval", startInterval)
	return time.Now(), 1 * time.Hour, nil
}

// func parseTime(alexaRequest models.AlexaRequest) (time.Duration, error) {
// 	duration := 1 * time.Hour
// 	if alexaRequest.Request.Intent.Name != "" {
// 		duration, err := time.Parse(alexaRequest.Request.Intent.Slots["StartTime"].Value, 1*time.Hour)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return duration, nil
// }

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

func parseDuration(alexaSlot models.AlexaSlot) time.Duration {
	if alexaSlot.Value == "" {
		return 1 * time.Hour
	}

	regex, err := regexp.Compile("PT(\\d+)(M|H)")
	if err != nil {
		return 1 * time.Hour
	}

	matches := regex.FindStringSubmatch(alexaSlot.Value)

	val, err := strconv.Atoi(matches[1])
	if err != nil {
		return 1 * time.Hour
	}

	switch matches[1] {
	case "M", "m":
		return time.Duration(val) * time.Minute
	case "H", "h":
		return time.Duration(val) * time.Hour
	}

	return 1 * time.Hour
}

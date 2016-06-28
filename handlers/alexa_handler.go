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

	alexaResp := models.NewAlexaResponse()

	fmt.Printf("request %#v\n", alexaRequest)
	startTime, duration, err := parse(alexaRequest)
	if err != nil {
		logger.Error("failed-parse-duration", err)
		alexaResp.OutputSpeech("I could not understand your request")
	} else {
		alexaResp.OutputSpeech(prepareResponse(startTime, duration))
		if err != nil {
			logger.Error("failed-prepare-response", err)
		}
	}

	json, err := alexaResp.String()
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.Write(json)

	logger.Info("request: " + fmt.Sprintf("%#v", alexaRequest))
}

func prepareResponse(startTime time.Time, duration time.Duration) string {
	return fmt.Sprintf("Hello world! %#v", startTime)
}

func parse(alexaRequest models.AlexaRequest) (time.Time, time.Duration, error) {
	startTime, err := parseTime(alexaRequest.Request.Intent.Slots["StartAt"])
	duration, err := parseDuration(alexaRequest.Request.Intent.Slots["Length"])
	return startTime, duration, err
}

func parseTime(alexaSlot models.AlexaSlot) (time.Time, error) {
	if alexaSlot.Value != "" {
		parsedTime, err := time.Parse("15:04", alexaSlot.Value)

		currentTime := time.Now()

		expectTime := time.Date(
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day(),
			parsedTime.Hour(),
			parsedTime.Minute(),
			0,
			0,
			currentTime.Location(),
		)

		if err != nil {
			return time.Now(), err
		}
		return expectTime, nil
	}
	return time.Now(), nil
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

func parseDuration(alexaSlot models.AlexaSlot) (time.Duration, error) {
	if alexaSlot.Value == "" {
		return 1 * time.Hour, nil
	}

	regex, err := regexp.Compile("PT(\\d+)(M|H)")
	if err != nil {
		return 1 * time.Hour, err
	}

	matches := regex.FindStringSubmatch(alexaSlot.Value)

	val, err := strconv.Atoi(matches[1])
	if err != nil {
		return 1 * time.Hour, err
	}

	switch matches[2] {
	case "M", "m":
		return time.Duration(val) * time.Minute, nil
	case "H", "h":
		return time.Duration(val) * time.Hour, nil
	}

	return 1 * time.Hour, nil
}

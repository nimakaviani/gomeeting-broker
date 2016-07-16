package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/nimakaviani/gomeeting-broker/models"
)

func PrepareResponse(room Room, startTime time.Time, duration time.Duration) string {
	hour, minute, _ := startTime.Clock()
	amPM := "AM"

	if hour > 12 {
		hour = hour - 12
		amPM = "PM"
	}

	startTimeString := fmt.Sprintf("%d:%d%s", hour, minute, amPM)
	if minute == 0 {
		startTimeString = fmt.Sprintf("%d%s", hour, amPM)
	}

	return fmt.Sprintf("Room %s is available from %s for %s", room.Name,
		startTimeString,
		humanizeLength(int(duration.Seconds())),
	)
}

func Parse(alexaRequest models.AlexaRequest) (time.Time, time.Duration, error) {
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

func humanizeLength(length int) string {
	switch {
	case length <= 60:
		return "1 minute"
	case length < 3600:
		return strconv.Itoa(length/60) + " minutes"
	case length > 3600:
		return strconv.Itoa(length/3600) + " hours"
	default:
		return "1 hour"
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

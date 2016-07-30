package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nimakaviani/gomeeting-broker/models"
)

const (
	FindRoom = "FindRoom"
	Locator  = "Locator"
)

func ParseAndGetResponse(alexaRequest models.AlexaRequest, config models.Config, datastore DataStore) (string, error) {

	switch alexaRequest.Request.Intent.Name {
	case FindRoom:
		startTime, err := parseTime(alexaRequest.Request.Intent.Slots["StartAt"], config)
		if err != nil {
			return "", err
		}

		duration, err := parseDuration(alexaRequest.Request.Intent.Slots["Length"])
		if err != nil {
			return "", err
		}

		calendar := NewGCalendar(GetClient(datastore))
		rooms := calendar.FindRoom(*startTime, duration)

		return prepareResponse(rooms[0], *startTime, duration), nil

	case Locator:
		roomName := strings.ToLower(alexaRequest.Request.Intent.Slots["RoomName"].Value)
		phrase, err := composeRoomLocation(roomName, config)
		if err != nil {
			return "", err
		}
		return phrase, nil
	}
	return "", fmt.Errorf("Could not find the defined slot")
}

func composeRoomLocation(roomName string, config models.Config) (string, error) {
	for _, roomObj := range config.Rooms {
		if roomObj.Name == roomName {
			return fmt.Sprintf(
				"%s is located on the %s floor, %s",
				roomObj.Name,
				roomObj.Floor,
				roomObj.Location), nil
		}
	}
	return "", fmt.Errorf("room not found")
}

func prepareResponse(room Room, startTime time.Time, duration time.Duration) string {
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

func parseTime(alexaSlot models.AlexaSlot, config models.Config) (*time.Time, error) {
	currentTime := time.Now()
	timeLocation, err := time.LoadLocation(config.Timezone)
	if err != nil {
		return nil, err
	}
	localTime := currentTime.In(timeLocation)

	if alexaSlot.Value != "" {
		parsedTime, err := time.Parse("15:04", alexaSlot.Value)

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
			return &localTime, err
		}
		return &expectTime, nil
	}
	return &localTime, nil
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

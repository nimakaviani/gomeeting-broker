package utils

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/calendar/v3"
)

type Room struct {
	Location  string
	Name      string
	Capacity  int
	Phone     string
	Extension string
}

type GCalendar interface {
	FindRoom(startTime time.Time, duration time.Duration) []Room
}

type gCalendar struct {
	client *http.Client
}

func NewGCalendar(client *http.Client) GCalendar {
	return &gCalendar{
		client: client,
	}
}

func (g *gCalendar) FindRoom(startTime time.Time, duration time.Duration) []Room {
	srv, err := calendar.New(g.client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	calList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Print("Unable to retrieve calendar Client %v", err)
	}

	calendarItems := g.filterCalendarRequestItems(calList.Items, "pivotal.io")

	freeBusy := calendar.NewFreebusyService(srv)
	freeBusyResponse, err := freeBusy.Query(&calendar.FreeBusyRequest{
		TimeMin: startTime.Format(time.RFC3339),
		TimeMax: startTime.Add(duration).Format(time.RFC3339),
		Items:   calendarItems,
	}).Do()

	matchedRooms := []Room{}
	for key, calendar := range freeBusyResponse.Calendars {
		if len(calendar.Busy) == 0 {
			matchedRoom := g.findRoom(calList, key)
			if matchedRoom != (Room{}) {
				matchedRooms = append(matchedRooms, matchedRoom)
			}
		}
	}

	return matchedRooms
}

func (g *gCalendar) filterCalendarRequestItems(items []*calendar.CalendarListEntry, filter string) []*calendar.FreeBusyRequestItem {
	regex, err := regexp.Compile(filter)
	if err != nil {
		log.Fatalf("Regex failed with error %#v", err)
	}

	calendarItems := []*calendar.FreeBusyRequestItem{}
	for _, item := range items {
		if !regex.Match([]byte(item.Id)) {
			continue
		}

		calendarItems = append(calendarItems, &calendar.FreeBusyRequestItem{Id: item.Id})
	}

	return calendarItems
}

func (g *gCalendar) findRoom(calendars *calendar.CalendarList, calendarKey string) Room {
	regex, err := regexp.Compile(`(.+) - ([a-zA-Z]+) .~(\d+) people\) (.+) ([a-zA-Z]+.*)`)
	if err != nil {
		log.Fatalf("Regex failed with error %#v", err)
	}

	for _, calendarItem := range calendars.Items {
		if calendarItem.Id == calendarKey {
			matchedStrings := regex.FindStringSubmatch(calendarItem.Summary)
			capacity, err := strconv.Atoi(matchedStrings[3])
			if err != nil {
				log.Fatalf("Regex failed with error %#v", err)
			}

			return Room{
				Location:  matchedStrings[1],
				Name:      matchedStrings[2],
				Capacity:  capacity,
				Phone:     matchedStrings[4],
				Extension: matchedStrings[5],
			}
		}
	}

	return Room{}
}

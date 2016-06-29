package utils

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"google.golang.org/api/calendar/v3"
)

type GCalendar interface {
	FindRoom(startTime time.Time, duration time.Duration)
}

type gCalendar struct {
	client *http.Client
}

func NewGCalendar(client *http.Client) GCalendar {
	return &gCalendar{
		client: client,
	}
}

func (g *gCalendar) FindRoom(startTime time.Time, duration time.Duration) {

	srv, err := calendar.New(g.client)

	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	calList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Print("Unable to retrieve calendar Client %v", err)
	}

	calendarItems := g.filterCalendarRequestItems(calList.Items, "gmail.com")

	freeBusy := calendar.NewFreebusyService(srv)
	freeBusyResponse, err := freeBusy.Query(&calendar.FreeBusyRequest{
		TimeMin: startTime.Format(time.RFC3339),
		TimeMax: startTime.Add(duration).Format(time.RFC3339),
		Items:   calendarItems,
	}).Do()

	fmt.Printf("Results %#v\n\n\n", freeBusyResponse)
}

func (g *gCalendar) filterCalendarRequestItems(items []*calendar.CalendarListEntry, filter string) []*calendar.FreeBusyRequestItem {
	regex, err := regexp.Compile(filter)
	if err != nil {
		log.Fatalf("Failed with error %#v", err)
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

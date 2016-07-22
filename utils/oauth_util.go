package utils

import (
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func GetClient(datastore DataStore) *http.Client {
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("Config field with error: %#v", err)
	}

	tok, err := datastore.LoadToken()
	if err != nil {
		log.Fatalf("reading from datastore failed: %v", err)
	}

	return config.Client(context.Background(), &tok)
}

func GetTokenURL(config *oauth2.Config) string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func GetConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	return config, nil
}

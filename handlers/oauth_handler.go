package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"google.golang.org/api/calendar/v3"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (h handler) OAuth(writer http.ResponseWriter, request *http.Request) {
	conf, err := getConfig()
	if err != nil {
		log.Fatal("failed")
	}
	url := getTokenURL(conf)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("<html><a href=\"" + url + "\" target=\"blank\">Click here to authorize access</a></html>"))
}

func (h handler) OAuthCallback(writer http.ResponseWriter, request *http.Request) {
	conf, err := getConfig()
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	query := request.URL.Query()
	state := query["state"][0]
	code := query["code"][0]

	if state != "" && code != "" {
		fmt.Printf("State: %s, Code: %s", state, code)
		tok := getToken(code, *conf)
		getClient(tok, conf)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Thank you! You are registered!"))
}

func getClient(tok *oauth2.Token, config *oauth2.Config) *http.Client {
	return config.Client(context.Background(), tok)
}

func getTokenURL(config *oauth2.Config) string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func getConfig() (*oauth2.Config, error) {
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

func getToken(code string, config oauth2.Config) *oauth2.Token {
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	fileName, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	saveToken(fileName, tok)
	return tok
}

func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), err
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

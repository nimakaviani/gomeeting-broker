package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func GetClient() *http.Client {
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("Config field with error: %#v", err)
	}

	cacheFile, err := GetTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file.\n First call <endpoint>/oauth %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		log.Fatalf("reading cached credentials file failed: %v", err)
	}
	return config.Client(context.Background(), tok)
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

func SaveToken(token *oauth2.Token) {
	file, err := GetTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}

	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetTokenCacheFile() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(workingDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), err
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

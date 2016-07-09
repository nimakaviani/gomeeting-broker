package utils

import (
	"encoding/json"
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

// needs some fixing here
type CalendarRunner struct {
	GCalendar GCalendar
}

func NewCalendarRunner() CalendarRunner {
	return CalendarRunner{}
}

func (r CalendarRunner) Run(singnals <-chan os.Signal, ready chan<- struct{}) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return err
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return err
	}
	r.GCalendar = NewGCalendar(getClient(ctx, config))
	close(ready)
	return nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file.\n First call <endpoint>/oauth %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		// log.Fatalf("reading cached credentials file failed: %v", err)
	}
	return config.Client(ctx, tok)
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(workingDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
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

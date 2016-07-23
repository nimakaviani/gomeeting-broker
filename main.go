package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"regexp"
	"syscall"

	"github.com/nimakaviani/gomeeting-broker/handlers"
	"github.com/nimakaviani/gomeeting-broker/models"
	"github.com/nimakaviani/gomeeting-broker/utils"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
)

var port = flag.String(
	"port",
	"8080",
	"Server port",
)

var sqlConnectionRegex = regexp.MustCompile(`^.+:\/\/(.+):(.+)@(.+):\d+\/(.+)$`)

type Credentials struct {
	URI string `json:"uri,omitempty"`
}

type ElephantSql struct {
	Credentials Credentials `json:"credentials,omitempty"`
}

type VCAPServices struct {
	Elephantsql []ElephantSql `json:"elephantsql,omitempty"`
}

type VCAPServicesStruct struct {
	VCAP_SERVICES VCAPServices `json:"VCAP_SERVICES"`
}

type DBCredentials struct {
	user     string
	password string
	host     string
	name     string
}

func main() {
	flag.Parse()
	httpServer := http_server.New(":"+*port, http.DefaultServeMux)

	config, err := models.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	dbCreds := getDBCreds()

	datastore := utils.NewDBDataStore(dbCreds.user, dbCreds.password, dbCreds.name, dbCreds.host)
	err = datastore.Init()
	if err != nil {
		panic(err)
	}
	defer datastore.Close()

	handler := handlers.NewHandler(config, datastore)
	http.HandleFunc("/findroom", handler.Alexa)
	http.HandleFunc("/google73d91fa1cfb6fa88.html", handler.Verify)
	http.HandleFunc("/oauth", handler.OAuth)
	http.HandleFunc("/oauthcallback", handler.OAuthCallback)
	http.HandleFunc("/favicon", handler.Icon)

	members := grouper.Members{
		{"httpserver", httpServer},
	}

	processes := grouper.NewOrdered(syscall.SIGINT, members)
	monitor := ifrit.Invoke(sigmon.New(processes))
	err = <-monitor.Wait()
	if err != nil {
		panic(err)
	}
}

func getDBCreds() DBCredentials {
	vcapServicesString := os.Getenv("VCAP_SERVICES")

	if vcapServicesString == "" {
		vcapServicesString = "{}"
	}

	vcapServices := VCAPServicesStruct{}

	err := json.Unmarshal([]byte(vcapServicesString), &vcapServices)
	if err != nil {
		panic(err)
	}

	if len(vcapServices.VCAP_SERVICES.Elephantsql) == 0 {
		return DBCredentials{
			user:     os.Getenv("DB_USER"),
			password: os.Getenv("DB_PASSWORD"),
			host:     os.Getenv("DB_HOST"),
			name:     os.Getenv("DB_NAME"),
		}
	}

	matches := sqlConnectionRegex.FindStringSubmatch(vcapServices.VCAP_SERVICES.Elephantsql[0].Credentials.URI)

	return DBCredentials{
		user:     matches[1],
		password: matches[2],
		host:     matches[3],
		name:     matches[4],
	}
}

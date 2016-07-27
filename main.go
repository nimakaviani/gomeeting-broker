package main

import (
	"flag"
	"net/http"
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

func main() {
	flag.Parse()

	config, err := models.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	dbCreds := utils.GetDBCreds()
	datastore := utils.NewDBDataStore(
		dbCreds.User,
		dbCreds.Password,
		dbCreds.Name,
		dbCreds.Host,
	)
	err = datastore.Init()
	if err != nil {
		panic(err)
	}
	defer datastore.Close()

	httpServer := http_server.New(":"+*port, http.DefaultServeMux)
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

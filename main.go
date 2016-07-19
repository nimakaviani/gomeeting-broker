package main

import (
	"flag"
	"net/http"
	"syscall"

	"github.com/nimakaviani/gomeeting-broker/handlers"
	"github.com/nimakaviani/gomeeting-broker/models"
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
	httpServer := http_server.New(":"+*port, http.DefaultServeMux)

	println("reading config ...")
	config, err := models.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	println("initializing handlers ...")
	handler := handlers.NewHandler(config)
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

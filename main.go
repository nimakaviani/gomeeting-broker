package main

import (
	"flag"
	"net/http"

	"github.com/nimakaviani/gomeeting-broker/handlers"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var port = flag.String(
	"port",
	"8080",
	"Server port",
)

func main() {
	flag.Parse()

	handler := handlers.NewHandler()
	http.HandleFunc("/", handler.Alexa)

	httpServer := http_server.New(":"+*port, http.DefaultServeMux)
	monitor := ifrit.Invoke(httpServer)
	err := <-monitor.Wait()
	if err != nil {
		panic(err)
	}
}

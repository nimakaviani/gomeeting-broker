package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nimakaviani/gomeeting-broker/utils"

	"golang.org/x/oauth2"
)

func (h handler) OAuth(writer http.ResponseWriter, request *http.Request) {
	conf, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("failed: %#v", err)
	}
	url := utils.GetTokenURL(conf)
	http.Redirect(writer, request, url, http.StatusFound)
}

func (h handler) OAuthCallback(writer http.ResponseWriter, request *http.Request) {
	conf, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	code := request.FormValue("code")
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	err = h.datastore.SaveToken(*token)
	if err != nil {
		log.Fatalf("Saving to datastore: %v", err)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Thank you! You are registered!"))
}

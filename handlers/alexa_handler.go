package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nimakaviani/gomeeting-broker/models"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Alexa(writer http.ResponseWriter, request *http.Request) {

	var alexaRequest models.AlexaRequest
	err := json.NewDecoder(request.Body).Decode(&alexaRequest)

	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("%v", alexaRequest))
}

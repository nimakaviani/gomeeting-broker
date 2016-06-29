package handlers

import "net/http"

func (h handler) Verify(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("google-site-verification: google73d91fa1cfb6fa88.html"))
}

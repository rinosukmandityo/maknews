package api

import (
	"log"
	"net/http"
)

func SetupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if _, e := w.Write(body); e != nil {
		log.Println(e)
	}
}

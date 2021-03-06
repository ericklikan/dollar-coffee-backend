package util

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func Respond(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.WithError(err).Warn()
	}
}

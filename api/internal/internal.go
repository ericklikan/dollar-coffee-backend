package internal

import (
	"errors"
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type internalSubrouter struct {
	util.CommonSubrouter
}

// This route is for internal uses only to update/get coffee, purchases etc

const prefix = "/internal"

func Setup(router *mux.Router, db *gorm.DB) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	internal := internalSubrouter{}
	internal.Router = router.
		PathPrefix(prefix).
		Subrouter()

	// Route to update and delete any coffees
	internal.Router.HandleFunc("/coffee", internal.coffeeHandler).Methods("POST", "DELETE")

	// Route to update amount paid on purchases and to view all purchases
	internal.Router.HandleFunc("/purchase", internal.purchaseHandler).Methods("GET", "PUT", "POST")
	return nil
}

func (sr *internalSubrouter) coffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalCoffeeHandler",
		"method":  r.Method,
	})
	logger.Warn("Unimplemented")
	util.Respond(w, util.Message("Unimplemented"))
}

func (sr *internalSubrouter) purchaseHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalPurchaseHandler",
		"method":  r.Method,
	})
	logger.Warn("Unimplemented")
	util.Respond(w, util.Message("Unimplemented"))
}

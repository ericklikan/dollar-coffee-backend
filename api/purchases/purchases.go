package purchases

import (
	"errors"
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/purchases"

func Setup(router *mux.Router, db *gorm.DB) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	subRouter := router.
		PathPrefix(prefix).
		Subrouter()

	// Set up auth middleware
	subRouter.Use(util.AuthMiddleware)

	// route for people to put in purchases, they should not be able to
	// put amount paid, this is done on internal route
	subRouter.HandleFunc("/purchase", PurchaseHandler).Methods("POST")
	return nil
}

func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "PurchaseHandler",
		"method":  r.Method,
	})
	logger.Warn("Unimplemented")
	util.Respond(w, util.Message("Unimplemented"))
}

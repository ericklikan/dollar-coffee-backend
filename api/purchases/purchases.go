package purchases

import (
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/purchases"

func Setup(router *mux.Router, db *gorm.DB) {
	subRouter := router.
		PathPrefix(prefix).
		Methods("GET", "POST").
		Subrouter()

	// Set up auth middlewate
	subRouter.Use(util.AuthMiddleware)

	// route for people to put in purchases, they should not be able to
	// put amount paid, this is done on internal route
	subRouter.HandleFunc("/purchase", PurchaseHandler).Methods("POST")
}

func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "PurchaseHandler",
		"method":  r.Method,
	})
	logger.Warn("Unimplemented")
	util.Respond(w, util.Message("Unimplemented"))
}

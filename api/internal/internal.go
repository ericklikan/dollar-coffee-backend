package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/ericklikan/dollar-coffee-backend/models"
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

	internal.Db = db
	internal.Router.Use(util.AuthMiddleware)

	// Route to update and delete any coffees
	internal.Router.HandleFunc("/coffee", internal.coffeeHandler).Methods("POST")

	// used to delete coffees from the menu
	internal.Router.HandleFunc("/coffee/{coffeeId}", internal.deleteCoffeeHandler).Methods("DELETE")

	// Route to update amount paid on purchases and to view all purchases
	internal.Router.HandleFunc("/purchase/{purchaseId}", internal.purchaseHandler).Methods("PATCH")
	return nil
}

func validateAdmin(ctx context.Context) bool {
	return ctx.Value("role") == "admin"
}

func (sr *internalSubrouter) coffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalCoffeeHandler",
		"method":  r.Method,
	})
	if !validateAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		util.Respond(w, util.Message("Invalid role type"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var coffeeInfo models.Coffee
	err := decoder.Decode(&coffeeInfo)
	if err != nil {
		logger.WithError(err).Warn()
		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message(err.Error()))
		return
	}

	err = coffeeInfo.Create(sr.Db)
	if err != nil {
		logger.WithError(err).Warn()
		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message(err.Error()))
		return
	}

	util.Respond(w, util.Message("Created new Coffee"))
}

func (sr *internalSubrouter) deleteCoffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalDeleteCoffeeHandler",
		"method":  r.Method,
	})

	if !validateAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		util.Respond(w, util.Message("Invalid role type"))
		return
	}

	logger.Warn("Unimplemented")
	w.WriteHeader(http.StatusNotImplemented)
	util.Respond(w, util.Message("Unimplemented"))
}

func (sr *internalSubrouter) purchaseHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalPurchaseHandler",
		"method":  r.Method,
	})

	if !validateAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		util.Respond(w, util.Message("Invalid role type"))
		return
	}

	logger.Warn("Unimplemented")
	w.WriteHeader(http.StatusNotImplemented)
	util.Respond(w, util.Message("Unimplemented"))
}

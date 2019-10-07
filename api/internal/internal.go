package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	// Route to update amount paid on purchases
	// Requires param: "amountPaid" in body
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

	vars := mux.Vars(r)
	requestedCoffee := vars["coffeeId"]
	coffeeId, err := strconv.Atoi(requestedCoffee)
	if err != nil {
		logger.Warn("Error parsing id")
		w.WriteHeader(http.StatusBadRequest)
		util.Respond(w, util.Message("Error parsing coffee id"))
		return
	}

	coffee := models.Coffee{}
	err = sr.Db.Table("coffees").Where("ID = ?", coffeeId).First(&coffee).Error
	if err != nil {
		logger.WithError(err).Warn("Database Error")

		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			util.Respond(w, util.Message("Not Found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message("Internal Error"))
		return
	}

	err = sr.Db.Delete(&coffee).Error
	if err != nil {
		logger.WithError(err).Warn("Error")
		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message("Internal Error"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	util.Respond(w, util.Message("Successfully deleted coffee"))
}

type PurchaseUpdateRequest struct {
	AmountPaid float32 `json:"amountPaid"`
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
	vars := mux.Vars(r)
	requestedPurchase := vars["purchaseId"]
	txId, err := strconv.Atoi(requestedPurchase)
	if err != nil {
		logger.Warn("Error parsing id")
		w.WriteHeader(http.StatusBadRequest)
		util.Respond(w, util.Message("Error parsing coffee id"))
		return
	}

	var reqData PurchaseUpdateRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqData)
	if err != nil {
		logger.WithError(err).Warn()
		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message(err.Error()))
		return
	}

	transaction := models.Transaction{}
	err = sr.Db.Table("transactions").Where("ID = ?", txId).First(&transaction).Error
	if err != nil {
		logger.WithError(err).Warn("Database Error")

		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			util.Respond(w, util.Message("Not Found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message("Internal Error"))
		return
	}

	transaction.AmountPaid = reqData.AmountPaid
	err = sr.Db.Save(&transaction).Error
	if err != nil {
		logger.WithError(err).Warn("Database Error")
		w.WriteHeader(http.StatusInternalServerError)
		util.Respond(w, util.Message("Internal Error"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	util.Respond(w, util.Message("Successfully updated purchase"))
}

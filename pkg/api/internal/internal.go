package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ericklikan/dollar-coffee-backend/pkg/api/util"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const pageSize = 10

type internalSubrouter struct {
	util.CommonSubrouter
}

type UpdateCoffeeRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float32 `json:"price"`
	InStock     *bool    `json:"inStock"`
}

type PurchaseUpdateRequest struct {
	AmountPaid float32 `json:"amountPaid"`
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
	internal.Router.HandleFunc("/coffee/{coffeeId}", internal.updateCoffeeHandler).Methods("PATCH", "DELETE")

	// Route to update amount paid on purchases
	// Requires param: "amountPaid" in body
	internal.Router.HandleFunc("/purchase/{purchaseId}", internal.purchaseHandler).Methods("PATCH")

	// Route to get information from all users
	internal.Router.HandleFunc("/users", internal.usersHandler).Methods("GET")
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
		util.Respond(w, http.StatusForbidden, util.Message("Invalid role type"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var coffeeInfo models.Coffee
	if err := decoder.Decode(&coffeeInfo); err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	// Validate info
	if len(coffeeInfo.Name) == 0 || coffeeInfo.Price == 0 {
		logger.Warn("Invalid coffee attributes")
		util.Respond(w, http.StatusBadRequest, util.Message("Invalid coffee attributes"))
		return
	}

	if err := sr.Db.Create(&coffeeInfo).Error; err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	util.Respond(w, http.StatusOK, util.Message("Created new Coffee"))
}

func (sr *internalSubrouter) updateCoffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalDeleteCoffeeHandler",
		"method":  r.Method,
	})

	if !validateAdmin(r.Context()) {
		util.Respond(w, http.StatusForbidden, util.Message("Invalid role type"))
		return
	}

	vars := mux.Vars(r)
	requestedCoffee := vars["coffeeId"]
	coffeeId, err := strconv.Atoi(requestedCoffee)
	if err != nil {
		logger.Warn("Error parsing id")
		util.Respond(w, http.StatusBadRequest, util.Message("Error parsing coffee id"))
		return
	}

	coffee := models.Coffee{}
	err = sr.Db.Table("coffees").
		Where("ID = ?", coffeeId).
		First(&coffee).Error

	if err != nil {
		logger.WithError(err).Warn("Database Error")

		if err == gorm.ErrRecordNotFound {
			util.Respond(w, http.StatusNotFound, util.Message("Not Found"))
			return
		}

		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	// Update coffee with new values
	if r.Method == "PATCH" {
		decoder := json.NewDecoder(r.Body)
		var newCoffeeInfo UpdateCoffeeRequest
		if err := decoder.Decode(&newCoffeeInfo); err != nil {
			logger.WithError(err).Warn()
			util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
			return
		}

		// Compare difference/default values
		if newCoffeeInfo.Name != nil {
			coffee.Name = *newCoffeeInfo.Name
		}
		if newCoffeeInfo.Price != nil {
			coffee.Price = *newCoffeeInfo.Price
		}
		if newCoffeeInfo.Description != nil {
			coffee.Description = *newCoffeeInfo.Description
		}
		if newCoffeeInfo.InStock != nil {
			coffee.InStock = *newCoffeeInfo.InStock
		}

		if err := sr.Db.Save(&coffee).Error; err != nil {
			logger.WithError(err).Warn("Error")
			util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
			return
		}
		util.Respond(w, http.StatusOK, util.Message("Successfully updated coffee"))
		return
	}

	// Delete the coffee
	if r.Method == "DELETE" {
		if err := sr.Db.Delete(&coffee).Error; err != nil {
			logger.WithError(err).Warn("Error")
			util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
			return
		}
		util.Respond(w, http.StatusOK, util.Message("Successfully deleted coffee"))
		return
	}
}

func (sr *internalSubrouter) purchaseHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalPurchaseHandler",
		"method":  r.Method,
	})

	if !validateAdmin(r.Context()) {
		util.Respond(w, http.StatusForbidden, util.Message("Invalid role type"))
		return
	}
	vars := mux.Vars(r)
	requestedPurchase := vars["purchaseId"]
	txId, err := strconv.Atoi(requestedPurchase)
	if err != nil {
		logger.Warn("Error parsing id")
		util.Respond(w, http.StatusBadRequest, util.Message("Error parsing coffee id"))
		return
	}

	var reqData PurchaseUpdateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	transaction := models.Transaction{}
	err = sr.Db.Table("transactions").Where("ID = ?", txId).First(&transaction).Error
	if err != nil {
		logger.WithError(err).Warn("Database Error")

		if err == gorm.ErrRecordNotFound {
			util.Respond(w, http.StatusNotFound, util.Message("Not Found"))
			return
		}

		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	transaction.AmountPaid = reqData.AmountPaid

	if err := sr.Db.Save(&transaction).Error; err != nil {
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	util.Respond(w, http.StatusOK, util.Message("Successfully updated purchase"))
}

func (sr *internalSubrouter) usersHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalUsersHandler",
		"method":  r.Method,
	})
	if !validateAdmin(r.Context()) {
		logger.Warn("Invalid role type")
		util.Respond(w, http.StatusForbidden, util.Message("Invalid role type"))
		return
	}

	offset := 0
	pageNumQuery := r.URL.Query().Get("page")
	if pageNum, err := strconv.Atoi(pageNumQuery); err == nil {
		offset = (pageNum - 1) * pageSize
	}

	users := make([]*models.User, 0, 10)
	err := sr.Db.
		Table("users").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).
		Error
	if err != nil {
		logger.WithError(err).Warn("Error retrieving values")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	if len(users) == 0 {
		logger.Warn("users not found")
		util.Respond(w, http.StatusNotFound, util.Message("Couldn't find any users"))
		return
	}

	response := util.Message("Users successfully queried")
	response["users"] = users

	util.Respond(w, http.StatusOK, response)
}

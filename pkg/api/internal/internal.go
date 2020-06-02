package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ericklikan/dollar-coffee-backend/pkg/api/util"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const pageSize = 10

type internalSubrouter struct {
	util.CommonSubrouter

	coffeeRepository   repository_interfaces.CoffeeRepository
	purchaseRepository repository_interfaces.TransactionsRepository
	userRepository     repository_interfaces.UserRepository
}

type UpdateCoffeeRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	InStock     *bool    `json:"inStock"`
}

type PurchaseUpdateRequest struct {
	AmountPaid float64 `json:"amountPaid"`
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

// This route is for internal uses only to update/get coffee, purchases etc

const prefix = "/internal"

func Setup(router *mux.Router, db *gorm.DB,
	coffeeRepository repository_interfaces.CoffeeRepository,
	purchaseRepository repository_interfaces.TransactionsRepository,
	userRepository repository_interfaces.UserRepository,
) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	internal := internalSubrouter{
		coffeeRepository:   coffeeRepository,
		userRepository:     userRepository,
		purchaseRepository: purchaseRepository,
	}
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

	// Route to update user role information
	internal.Router.HandleFunc("/users/{userId}/role", internal.updateUserRoleHandler).Methods("PATCH")

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

	tx := sr.Db.Begin()
	if err := sr.coffeeRepository.CreateCoffee(tx, &coffeeInfo); err != nil {
		tx.Rollback()
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}
	tx.Commit()

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

	tx := sr.Db.Begin()
	coffeeMap, err := sr.coffeeRepository.GetCoffeesByIds(tx, []string{requestedCoffee})
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	if _, doesCoffeeExist := coffeeMap[requestedCoffee]; !doesCoffeeExist {
		tx.Rollback()
		util.Respond(w, http.StatusNotFound, util.Message("Couldn't find coffee"))
		return
	}

	coffee := coffeeMap[requestedCoffee]

	// Update coffee with new values
	if r.Method == "PATCH" {
		decoder := json.NewDecoder(r.Body)
		var newCoffeeInfo UpdateCoffeeRequest
		if err := decoder.Decode(&newCoffeeInfo); err != nil {
			tx.Rollback()
			logger.WithError(err).Warn("Error decoding JSON")
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

		if err := sr.coffeeRepository.UpdateCoffee(tx, coffee); err != nil {
			tx.Rollback()
			logger.WithError(err).Warn("Error")
			util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
			return
		}
		tx.Commit()
		util.Respond(w, http.StatusOK, util.Message("Successfully updated coffee"))
		return
	}

	// Delete the coffee
	if r.Method == "DELETE" {
		if err := sr.coffeeRepository.DeleteCoffee(tx, requestedCoffee); err != nil {
			tx.Rollback()
			logger.WithError(err).Warn("Error")
			util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
			return
		}
		tx.Commit()
		util.Respond(w, http.StatusOK, util.Message("Successfully deleted coffee"))
		return
	}

	tx.Rollback()
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

	var reqData PurchaseUpdateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	tx := sr.Db.Begin()
	transactionsMap, err := sr.purchaseRepository.GetTransactionsByIds(tx, []string{requestedPurchase})
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	var transaction *models.Transaction
	var doesTxExist bool
	if transaction, doesTxExist = transactionsMap[requestedPurchase]; !doesTxExist {
		tx.Rollback()
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusNotFound, util.Message("Transaction not found"))
		return
	}

	transaction.AmountPaid = reqData.AmountPaid

	if err := sr.purchaseRepository.UpdateTransaction(tx, transaction); err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}
	tx.Commit()

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

	pageNum := 0
	pageNumQuery := r.URL.Query().Get("page")
	if pageNumInt, err := strconv.Atoi(pageNumQuery); err == nil {
		pageNum = pageNumInt - 1
	}

	query := repository_interfaces.UsersPageQuery{}
	query.Page = pageNum
	query.PageSize = pageSize

	if role := r.URL.Query().Get("role"); role == "admin" || role == "user" {
		query.Role = &role
	}

	tx := sr.Db.Begin()
	users, err := sr.userRepository.GetUsersPaginated(tx, &query)
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
	response["pageSize"] = len(users)

	util.Respond(w, http.StatusOK, response)
}

func (sr *internalSubrouter) updateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "InternalUpdateUserRole",
		"method":  r.Method,
	})

	if !validateAdmin(r.Context()) {
		util.Respond(w, http.StatusForbidden, util.Message("Invalid role type"))
		return
	}
	vars := mux.Vars(r)
	requestedUser := vars["userId"]

	var reqData UpdateRoleRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	if reqData.Role != "admin" && reqData.Role != "user" {
		logger.Warn("Invalid Role")
		util.Respond(w, http.StatusBadRequest, util.Message("Invalid Role"))
		return
	}

	// retrieve user
	tx := sr.Db.Begin()
	usersMap, err := sr.userRepository.GetUsersByIds(tx, []string{requestedUser})
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	var user *models.User
	var doesUserExist bool
	if user, doesUserExist = usersMap[requestedUser]; !doesUserExist {
		tx.Rollback()
		logger.WithError(err).Warn("user not found")
		util.Respond(w, http.StatusNotFound, util.Message("User not found"))
		return
	}

	user.Role = reqData.Role
	if err := sr.userRepository.UpdateUser(tx, user); err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Database Error")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	tx.Commit()
	util.Respond(w, http.StatusOK, util.Message("Successfully updated user role"))
}

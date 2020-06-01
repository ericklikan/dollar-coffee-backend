package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ericklikan/dollar-coffee-backend/pkg/api/util"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/auth"

type authSubrouter struct {
	util.CommonSubrouter

	userRepository repository_interfaces.UserRepository
}

func Setup(router *mux.Router, db *gorm.DB, userRepository repository_interfaces.UserRepository) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	auth := authSubrouter{
		userRepository: userRepository,
	}
	auth.Router = router.
		PathPrefix(prefix).
		Subrouter()
	auth.Db = db

	auth.Router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	auth.Router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	return nil
}

func (sr *authSubrouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "LoginHandler",
		"method":  r.Method,
	})

	// Login requires 2 pieces of data:
	// - email
	// - password
	decoder := json.NewDecoder(r.Body)
	var userInfo models.User
	err := decoder.Decode(&userInfo)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}
	if len(userInfo.Email) == 0 || len(userInfo.Password) == 0 {
		logger.Warn("Missing email or password")
		util.Respond(w, http.StatusBadRequest, util.Message("missing user"))
		return
	}

	err = userInfo.Login(sr.Db)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusUnauthorized, util.Message(err.Error()))
		return
	}

	response := util.Message(fmt.Sprintf("Logged In as %s", userInfo.FirstName))
	response["token"] = userInfo.Token

	util.Respond(w, http.StatusOK, response)
}

func (sr *authSubrouter) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "RegisterHandler",
		"method":  r.Method,
	})

	// Register requires 4 pieces of data:
	// - firstName
	// - lastName
	// - email
	// - password
	// - phone (OPTIONAL)
	decoder := json.NewDecoder(r.Body)
	var userInfo models.User
	err := decoder.Decode(&userInfo)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	// make email to lower
	userInfo.Email = strings.ToLower(userInfo.Email)

	if len(userInfo.FirstName) == 0 || len(userInfo.LastName) == 0 {
		logger.Warn("you must have a first name and a last name")
		util.Respond(w, http.StatusBadRequest, util.Message("Invalid first or last name"))
		return
	}

	if err := userInfo.Create(sr.Db); err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	response := util.Message("Created User")
	response["token"] = userInfo.Token

	util.Respond(w, http.StatusCreated, response)
}

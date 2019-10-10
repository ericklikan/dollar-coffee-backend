package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/ericklikan/dollar-coffee-backend/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/auth"

type authSubrouter struct {
	util.CommonSubrouter
}

func Setup(router *mux.Router, db *gorm.DB) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	auth := authSubrouter{}
	auth.Router = router.
		PathPrefix(prefix).
		Subrouter()
	auth.Db = db

	auth.Router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	auth.Router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")

	// TODO: Implement refresh tokens to get new jwt
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

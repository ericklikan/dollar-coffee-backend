package auth

import (
	"errors"
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
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
	return nil
}

func (s *authSubrouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "LoginHandler",
		"method":  r.Method,
	})

	logger.Warn("Unimplemented")
}

func (s *authSubrouter) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "RegisterHandler",
		"method":  r.Method,
	})

	logger.Warn("Unimplemented")
}

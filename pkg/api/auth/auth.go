package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/util"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

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

	tx := sr.Db.Begin()
	user, err := sr.userRepository.GetUserByEmail(tx, userInfo.Email)
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn()
		if err == gorm.ErrRecordNotFound {
			util.Respond(w, http.StatusNotFound, util.Message(err.Error()))
			return
		}
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}
	tx.Commit()

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		util.Respond(w, http.StatusUnauthorized, util.Message("Could not verify user password"))
		return
	}

	// Queried user is now valid
	user.Password = ""

	//Create JWT token
	tk := &models.Token{
		UserId: user.ID,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	user.Token = tokenString //Store the token in the response

	response := util.Message(fmt.Sprintf("Logged In as %s", user.FirstName))
	response["token"] = user.Token
	response["userId"] = user.ID

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
	var userInfo *models.User
	err := decoder.Decode(&userInfo)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	// make email to lower
	userInfo.Email = strings.ToLower(userInfo.Email)

	if err := userInfo.Validate(); err != nil || len(userInfo.FirstName) == 0 || len(userInfo.LastName) == 0 {
		logger.Warn("you must have a first name and a last name")
		util.Respond(w, http.StatusBadRequest, util.Message("Invalid first or last name"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
	}
	userInfo.Password = string(hashedPassword)

	// Generate UUID
	userInfo.ID = uuid.New()

	tx := sr.Db.Begin()
	if err := sr.userRepository.CreateUser(tx, userInfo); err != nil {
		tx.Rollback()
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}
	tx.Commit()

	//Create new JWT token for the newly registered account and default to role type as user
	tk := &models.Token{
		UserId: userInfo.ID,
		Role:   "user",
	}

	// HS256 is a symmetric key encryption algorithm. The same token password that is used to sign the token is used to verify the token
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	userInfo.Token = tokenString

	userInfo.Password = "" //delete password

	response := util.Message("Created User")
	response["token"] = userInfo.Token
	response["userId"] = userInfo.ID

	util.Respond(w, http.StatusCreated, response)
}

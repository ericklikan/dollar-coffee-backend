package menu

import (
	"net/http"
	"strconv"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/menu"
const pageSize = 10

type menuSubrouter struct {
	util.CommonSubrouter
}

func Setup(router *mux.Router, db *gorm.DB) error {
	subRouter := menuSubrouter{}
	subRouter.Router = router.PathPrefix(prefix).Subrouter()
	subRouter.Db = db

	// Get all the coffees that are available
	// TODO: refactor page number to page token using coffee ids
	subRouter.Router.HandleFunc("/coffee", subRouter.CoffeeHandler).Methods("GET")
	return nil
}

type CoffeeResponse struct {
	ID          uint
	Name        string
	Description string
	Price       float32
	InStock     bool
}

func (router *menuSubrouter) CoffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "CoffeeHandler",
		"method":  r.Method,
	})

	// Find page offset
	offset := 0
	pageNumQuery := r.URL.Query().Get("page")
	if pageNum, err := strconv.Atoi(pageNumQuery); err == nil {
		offset = (pageNum - 1) * pageSize
	}

	// query coffees
	coffees := []CoffeeResponse{}
	err := router.Db.Table("coffees").
		Select([]string{"ID", "name", "description", "price", "in_stock"}).
		Offset(offset).
		Limit(pageSize).
		Find(&coffees).
		Error
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message("Error getting coffees"))
	}
	if len(coffees) == 0 {
		logger.Warn("coffees not found")
		util.Respond(w, http.StatusNotFound, util.Message("Couldn't find coffees"))
		return
	}

	util.Respond(w, http.StatusOK, map[string]interface{}{
		"coffees": coffees,
	})
}

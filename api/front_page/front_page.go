package front_page

import (
	"net/http"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/ericklikan/dollar-coffee-backend/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/menu"

type frontPageRouter struct {
	util.CommonSubrouter
}

func Setup(router *mux.Router, db *gorm.DB) error {
	subRouter := frontPageRouter{}
	subRouter.Router = router.PathPrefix(prefix).Subrouter()
	subRouter.Db = db

	subRouter.Router.HandleFunc("/coffee", subRouter.CoffeeHandler).Methods("GET")
	return nil
}

func (router *frontPageRouter) CoffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "PurchaseHandler",
		"method":  r.Method,
	})

	coffees := []models.Coffee{}
	dbTxInfo := router.Db.Table("coffees").Limit(10).Find(&coffees)
	if dbTxInfo.Error != nil {
		logger.WithError(dbTxInfo.Error).Warn()
		util.Respond(w, util.Message("Error getting coffees"))
	}

	util.Respond(w, map[string]interface{}{
		"coffees": coffees,
	})
}

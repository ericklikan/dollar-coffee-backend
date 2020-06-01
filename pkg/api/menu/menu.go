package menu

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericklikan/dollar-coffee-backend/pkg/api/util"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/menu"
const pageSize = 50

type MenuSubrouter struct {
	util.CommonSubrouter

	coffeeRepository repository_interfaces.CoffeeRepository
}

type CoffeeResponse struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	InStock     bool
	UpdatedAt   time.Time `json:"-"`
}

func Setup(router *mux.Router, db *gorm.DB, coffeeRepository repository_interfaces.CoffeeRepository) error {
	subRouter := MenuSubrouter{
		coffeeRepository: coffeeRepository,
	}
	subRouter.Router = router.PathPrefix(prefix).Subrouter()
	subRouter.Db = db

	// Get all the coffees that are available
	subRouter.Router.HandleFunc("", subRouter.CoffeeHandler).Methods("GET")
	return nil
}

func (router *MenuSubrouter) CoffeeHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "CoffeeHandler",
		"method":  r.Method,
	})

	var err error

	// Find pagenum
	pageNum := 0
	pageNumQuery := r.URL.Query().Get("page")
	if pageNumInt, err := strconv.Atoi(pageNumQuery); err == nil {
		pageNum = pageNumInt - 1
	}

	query := repository_interfaces.CoffeePageQuery{}
	query.PageSize = pageSize
	query.Page = pageNum

	// find in stock coffees only
	inStockQuery := r.URL.Query().Get("in_stock")
	if inStockBool, err := strconv.ParseBool(inStockQuery); err == nil {
		query.InStock = &inStockBool
	}

	// query coffees
	tx := router.Db.Begin()
	coffees, err := router.coffeeRepository.GetCoffeesPaginated(tx, &query)
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message("Error getting coffees"))
		return
	}
	tx.Commit()

	if len(coffees) == 0 {
		logger.Warn("coffees not found")
		util.Respond(w, http.StatusNotFound, util.Message("Couldn't find coffees"))
		return
	}

	res := make([]*CoffeeResponse, 0, len(coffees))
	for _, coffee := range coffees {
		res = append(res, &CoffeeResponse{
			ID:          coffee.ID,
			Name:        coffee.Name,
			Description: coffee.Description,
			Price:       coffee.Price,
			InStock:     coffee.InStock,
		})
	}

	util.Respond(w, http.StatusOK, map[string]interface{}{
		"coffees":   res,
		"page_size": pageSize,
	})
}

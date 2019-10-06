package api

import (
	"github.com/ericklikan/dollar-coffee-backend/api/auth"
	"github.com/ericklikan/dollar-coffee-backend/api/front_page"
	"github.com/ericklikan/dollar-coffee-backend/api/purchases"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Router *mux.Router
}

func NewServer(router *mux.Router, db *gorm.DB) error {
	server := Server{
		Router: router,
	}

	err := front_page.Setup(server.Router, db)
	if err != nil {
		return err
	}

	purchases.Setup(server.Router, db)

	err = auth.Setup(server.Router, db)
	if err != nil {
		return err
	}
	return nil
}

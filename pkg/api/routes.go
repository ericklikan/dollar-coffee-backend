package api

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/auth"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/internal"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/menu"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/purchases"

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

	err := menu.Setup(server.Router, db)
	if err != nil {
		return err
	}

	err = purchases.Setup(server.Router, db)
	if err != nil {
		return err
	}

	err = auth.Setup(server.Router, db)
	if err != nil {
		return err
	}

	err = internal.Setup(server.Router, db)
	if err != nil {
		return err
	}

	return nil
}

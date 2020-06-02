package api

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/auth"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/internal"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/menu"
	"github.com/ericklikan/dollar-coffee-backend/pkg/api/purchases"
	repository "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/impl"
	"github.com/go-redis/redis/v7"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Router *mux.Router
}

func NewServer(router *mux.Router, db *gorm.DB, redis *redis.Client) error {
	server := Server{
		Router: router,
	}

	// Repository setups
	coffeeRepository := repository.NewCoffeeRepository(db, redis)
	transactionRepository := repository.NewTransactionsRepository(db)
	userRepository := repository.NewUserRepository(db)

	// module setups
	err := menu.Setup(server.Router, db, coffeeRepository)
	if err != nil {
		return err
	}

	err = purchases.Setup(server.Router, db, coffeeRepository, transactionRepository)
	if err != nil {
		return err
	}

	err = auth.Setup(server.Router, db, userRepository)
	if err != nil {
		return err
	}

	err = internal.Setup(server.Router, db, coffeeRepository, transactionRepository, userRepository)
	if err != nil {
		return err
	}

	return nil
}

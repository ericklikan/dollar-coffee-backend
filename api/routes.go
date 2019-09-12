package api

import (
	"fmt"

	frontPage "github.com/ericklikan/dollar-coffee-backend/api/front_page"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer(router *mux.Router) (*Server, error) {
	server := Server{
		Router: router,
	}

	err := frontPage.Setup("/front_page", server.Router)
	if err != nil {
		//TODO replace with log
		fmt.Println("Error setting up front page route")
		return nil, err
	}
	return &server, nil
}

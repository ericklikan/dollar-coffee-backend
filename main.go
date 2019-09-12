package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	routes "github.com/ericklikan/dollar-coffee-backend/api"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("running in debug mode: port 5000")
		port = "5000"
	}

	router := mux.NewRouter()
	_, err := routes.NewServer(router)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

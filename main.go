package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	routes "github.com/ericklikan/dollar-coffee-backend/api"
	"github.com/gorilla/handlers"
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
	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

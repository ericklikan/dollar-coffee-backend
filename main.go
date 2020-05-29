package main

import (
	"fmt"
	"net/http"
	"os"

	routes "github.com/ericklikan/dollar-coffee-backend/pkg/api"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup environment
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Warn(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Info("running in debug mode: port 5000")
		port = "5000"
	}

	db, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	err = routes.NewServer(router, db)
	if err != nil {
		log.Fatal(err)
	}

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Infof("Started server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func setupDatabase() (*gorm.DB, error) {
	dbUri := os.Getenv("DATABASE_URL")

	log.Infof("Connecting to postgresdb: %s", dbUri)
	dbConn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		return nil, err
	}

	// Database migrations
	dbConn = dbConn.AutoMigrate(
		models.Coffee{},
		models.PurchaseItem{},
		models.Transaction{},
		models.User{},
	)
	if dbConn.Error != nil {
		return nil, dbConn.Error
	}

	dbConn = dbConn.Model(&models.PurchaseItem{}).AddForeignKey("transaction_id", "transactions(id)", "RESTRICT", "RESTRICT")
	if dbConn.Error != nil {
		// Will error if foreign key is already set up
		log.WithError(dbConn.Error).Warn()
	}

	dbConn = dbConn.Model(&models.Transaction{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	if dbConn.Error != nil {
		// Will error if foreign key is already set up
		log.WithError(dbConn.Error).Warn()
	}

	return dbConn, nil
}

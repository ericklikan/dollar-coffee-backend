package purchases

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ericklikan/dollar-coffee-backend/pkg/api/util"
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/purchases"
const pageSize = 10

type PurchaseSubRouter struct {
	util.CommonSubrouter

	coffeeRepository   repository_interfaces.CoffeeRepository
	purchaseRepository repository_interfaces.TransactionsRepository
}

type PurchaseItem struct {
	CoffeeId      uint   `json:"coffeeId"`
	CoffeeOptions string `json:"options"`
}

// Requests

type PurchaseRequest struct {
	Coffees []PurchaseItem `json:"items"`
}

// Responses

type PurchaseHistoryResponse struct {
	ID            uint                   `json:"transactionId"`
	AmountPaid    float64                `json:"amountPaid"`
	Total         float64                `json:"total"`
	CreatedAt     time.Time              `json:"purchaseDate"`
	PurchaseItems []*models.PurchaseItem `json:"items"`
}

func Setup(router *mux.Router, db *gorm.DB, coffeeRepository repository_interfaces.CoffeeRepository, transactionRepository repository_interfaces.TransactionsRepository) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	purchase := PurchaseSubRouter{
		coffeeRepository:   coffeeRepository,
		purchaseRepository: transactionRepository,
	}
	purchase.Router = router.
		PathPrefix(prefix).
		Subrouter()

	purchase.Db = db
	// Set up auth middleware
	purchase.Router.Use(util.AuthMiddleware)

	// route for people to put in purchases, they should not be able to
	// put amount paid, this is done on internal route
	purchase.Router.HandleFunc("/purchase", purchase.PurchaseHandler).Methods("POST")

	// /purchases/{userId} will get the purchase history for that user.
	// query parameters should be pageNum
	// TODO: refactor page number to page token using purchase ids
	purchase.Router.HandleFunc("/user/{userId}", purchase.PurchaseHistoryHandler).Methods("GET")
	return nil
}

func (sr *PurchaseSubRouter) PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "PurchaseHandler",
		"method":  r.Method,
	})

	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		logger.Warn("Error parsing uuid")
		util.Respond(w, http.StatusInternalServerError, util.Message("Error parsing user id header"))
		return
	}

	var reqData PurchaseRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	if len(reqData.Coffees) == 0 {
		logger.Warn("Order items can't be empty")
		util.Respond(w, http.StatusBadRequest, util.Message("Invalid request, order can't be empty"))
		return
	}

	coffeeIdsMap := make(map[string]bool)

	purchaseItems := make([]*models.PurchaseItem, 0, len(reqData.Coffees))
	for _, item := range reqData.Coffees {
		coffeeIdsMap[strconv.FormatUint(uint64(item.CoffeeId), 10)] = true
		purchaseItems = append(purchaseItems, &models.PurchaseItem{
			CoffeeId:   item.CoffeeId,
			TypeOption: item.CoffeeOptions,
		})
	}

	coffeeIds := make([]string, 0, len(coffeeIdsMap))
	for coffeeId := range coffeeIdsMap {
		coffeeIds = append(coffeeIds, coffeeId)
	}

	tx := sr.Db.Begin()
	coffeesMap, err := sr.coffeeRepository.GetCoffeesByIds(tx, coffeeIds)
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Error retrieving coffees")
		util.Respond(w, http.StatusInternalServerError, util.Message("InternalError"))
		return
	}

	totalPrice := 0.0
	for _, purchaseItem := range purchaseItems {
		coffee, exists := coffeesMap[strconv.FormatUint(uint64(purchaseItem.CoffeeId), 10)]
		if !exists {
			tx.Rollback()
			logger.Warn("Error coffee doesn't exist")
			util.Respond(w, http.StatusInternalServerError, util.Message("InternalError"))
			return
		}
		purchaseItem.Price = coffee.Price
		totalPrice += purchaseItem.Price
	}

	purchase := models.Transaction{
		UserId: userId,
		Items:  purchaseItems,
		Total:  totalPrice,
	}

	err = sr.purchaseRepository.CreateTransaction(tx, &purchase)
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

	tx.Commit()
	util.Respond(w, http.StatusOK, util.Message("Purchase Confirmed"))
}

func (sr *PurchaseSubRouter) PurchaseHistoryHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"request": "PurchaseHistoryHandler",
		"method":  r.Method,
	})
	vars := mux.Vars(r)
	requestedUserId := vars["userId"]

	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		logger.Warn("Error parsing uuid")
		util.Respond(w, http.StatusInternalServerError, util.Message("Error parsing user id path"))
		return
	}
	role, ok := r.Context().Value("role").(string)
	if !ok {
		logger.Warn("Error parsing role")
		util.Respond(w, http.StatusInternalServerError, util.Message("Invalid role"))
		return
	}
	if role != "admin" && userId.String() != requestedUserId {
		logger.Warn("Unauthorized user")
		util.Respond(w, http.StatusUnauthorized, util.Message("You can't view this information"))
		return
	}

	pageNum := 0
	pageNumQuery := r.URL.Query().Get("page")
	if pageNumInt, err := strconv.Atoi(pageNumQuery); err == nil {
		pageNum = pageNumInt - 1
	}

	sortKey := "created_at"
	sortDirection := "DESC"
	query := repository_interfaces.PurchasePageQuery{
		UserId:        &requestedUserId,
		Sort:          &sortKey,
		SortDirection: &sortDirection,
	}
	query.Page = pageNum
	query.PageSize = pageSize

	tx := sr.Db.Begin()
	dbPurchases, err := sr.purchaseRepository.GetTransactionsPaginated(tx, &query)
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Warn("Error retrieving values")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}
	tx.Commit()

	if len(dbPurchases) == 0 {
		logger.Warn("purchases not found")
		util.Respond(w, http.StatusNotFound, util.Message("Couldn't find any purchases"))
		return
	}

	purchases := make([]*PurchaseHistoryResponse, 0, len(dbPurchases))
	for _, purchase := range dbPurchases {
		purchaseItem := PurchaseHistoryResponse{
			ID:            purchase.ID,
			AmountPaid:    purchase.AmountPaid,
			Total:         purchase.Total,
			CreatedAt:     purchase.CreatedAt,
			PurchaseItems: purchase.Items,
		}
		purchases = append(purchases, &purchaseItem)
	}

	response := util.Message("Purchases successfully queried")
	response["purchases"] = purchases
	response["page_size"] = pageSize

	util.Respond(w, http.StatusOK, response)
}

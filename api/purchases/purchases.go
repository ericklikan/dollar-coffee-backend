package purchases

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ericklikan/dollar-coffee-backend/api/util"
	"github.com/ericklikan/dollar-coffee-backend/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/purchases"
const pageSize = 10

type PurchaseSubRouter struct {
	util.CommonSubrouter
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
	ID            uint                  `json:"transactionId"`
	AmountPaid    float32               `json:"amountPaid"`
	CreatedAt     time.Time             `json:"purchaseDate"`
	PurchaseItems []models.PurchaseItem `json:"items"`
}

func Setup(router *mux.Router, db *gorm.DB) error {
	if db == nil || router == nil {
		err := errors.New("db or router is nil")
		log.WithError(err).Warn()
		return err
	}

	purchase := PurchaseSubRouter{}
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

	purchaseItems := make([]models.PurchaseItem, 0, len(reqData.Coffees))
	for _, item := range reqData.Coffees {
		purchaseItems = append(purchaseItems, models.PurchaseItem{
			CoffeeId:   item.CoffeeId,
			TypeOption: item.CoffeeOptions,
		})
	}

	purchase := models.Transaction{
		UserId: userId,
		Items:  purchaseItems,
	}

	err = sr.Db.Create(&purchase).Error
	if err != nil {
		logger.WithError(err).Warn()
		util.Respond(w, http.StatusInternalServerError, util.Message(err.Error()))
		return
	}

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

	offset := 0
	pageNumQuery := r.URL.Query().Get("page")
	if pageNum, err := strconv.Atoi(pageNumQuery); err == nil {
		offset = (pageNum - 1) * pageSize
	}

	purchases := make([]PurchaseHistoryResponse, 0)
	err := sr.Db.
		Table("transactions").
		Select([]string{"id", "amount_paid", "created_at"}).
		Where("user_id = ?", requestedUserId).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&purchases).
		Error
	if err != nil {
		logger.WithError(err).Warn("Error retrieving values")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	if len(purchases) == 0 {
		logger.Warn("purchases not found")
		util.Respond(w, http.StatusNotFound, util.Message("Couldn't find any purchases"))
		return
	}

	purchaseIds := make([]uint, 0, len(purchases))
	for _, purchase := range purchases {
		purchaseIds = append(purchaseIds, purchase.ID)
	}

	purchaseItems := []models.PurchaseItem{}
	err = sr.Db.
		Table("purchase_items").
		Where("transaction_id in (?)", purchaseIds).
		Find(&purchaseItems).
		Error
	if err != nil {
		logger.WithError(err).Warn("Error retrieving values")
		util.Respond(w, http.StatusInternalServerError, util.Message("Internal Error"))
		return
	}

	// Put it in a map to reduce queries
	purchaseItemMap := make(map[uint][]models.PurchaseItem)
	for _, purchaseItem := range purchaseItems {
		purchaseItemMap[purchaseItem.TransactionId] = append(purchaseItemMap[purchaseItem.TransactionId], purchaseItem)
	}

	for i, purchase := range purchases {
		purchases[i].PurchaseItems = purchaseItemMap[purchase.ID]
	}

	response := util.Message("Purchases successfully queried")
	response["purchases"] = purchases

	util.Respond(w, http.StatusOK, response)
}

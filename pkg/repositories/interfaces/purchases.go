package repository_interfaces

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

type PurchasePageQuery struct {
	PageQuery

	UserId        *string
	Sort          *string
	SortDirection *string //If you want to query with sort direction you need a sort key
}

type TransactionsRepository interface {
	CreateTransaction(tx *gorm.DB, transaction models.Transaction) error
	GetTransactionsPaginated(tx *gorm.DB, query PurchasePageQuery) ([]*models.Transaction, error)
	UpdateTransaction(tx *gorm.DB, transaction *models.Transaction) error
	DeleteTransaction(tx *gorm.DB, transactionId string) error
}

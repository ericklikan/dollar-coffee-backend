package repository

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/ericklikan/dollar-coffee-backend/pkg/persistence"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/jinzhu/gorm"
)

type TransactionsRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionsRepository(db *gorm.DB) repository_interfaces.TransactionsRepository {
	return &TransactionsRepositoryImpl{
		db: db,
	}
}

func (repo *TransactionsRepositoryImpl) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	return persistence.CreateTransaction(tx, transaction)
}

func (repo *TransactionsRepositoryImpl) GetTransactionsByIds(tx *gorm.DB, transactionIds []string) (map[string]*models.Transaction, error) {
	return persistence.GetTransactionsByID(tx, transactionIds)
}

func (repo *TransactionsRepositoryImpl) GetTransactionsPaginated(tx *gorm.DB, query *repository_interfaces.PurchasePageQuery) ([]*models.Transaction, error) {
	return persistence.GetTransactionsPaginated(tx, query.PageSize, query.Page, query.UserId)
}

func (repo *TransactionsRepositoryImpl) UpdateTransaction(tx *gorm.DB, purchase *models.Transaction) error {
	return persistence.UpdateTransaction(tx, purchase)
}

func (repo *TransactionsRepositoryImpl) DeleteTransaction(tx *gorm.DB, purchaseId string) error {
	return persistence.DeleteTransaction(tx, purchaseId)
}

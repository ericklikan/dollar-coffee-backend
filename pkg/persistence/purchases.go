package persistence

import (
	"fmt"
	"strconv"

	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

func CreateTransaction(tx *gorm.DB, purchase *models.Transaction) error {
	return tx.Create(purchase).Error
}

func GetTransactionsByID(tx *gorm.DB, purchaseIds []string) (map[string]*models.Transaction, error) {
	var purchases []*models.Transaction
	if err := tx.
		Where("id in (?)", purchaseIds).
		Find(&purchases).Error; err != nil {
		return nil, err
	}
	purchaseMap := make(map[string]*models.Transaction)
	for _, purchase := range purchases {
		purchaseMap[strconv.FormatUint(uint64(purchase.ID), 10)] = purchase
	}
	return purchaseMap, nil
}

func GetTransactionsPaginated(tx *gorm.DB, pageSize int, page int, userId *string, sortKey *string, sortDirection *string) ([]*models.Transaction, error) {
	var purchases []*models.Transaction
	q := tx.
		Offset(page * pageSize).
		Limit(pageSize).
		Preload("Items")

	if userId != nil {
		q = q.Where("user_id = ?", *userId)
	}

	if sortKey != nil && sortDirection != nil {
		q = q.Order(fmt.Sprintf("%s %s", *sortKey, *sortDirection))
	} else {
		q = q.Order("created_at DESC")
	}

	if err := q.Find(&purchases).Error; err != nil {
		return nil, err
	}

	return purchases, nil
}

func UpdateTransaction(tx *gorm.DB, purchase *models.Transaction) error {
	return tx.Save(purchase).Error
}

func DeleteTransaction(tx *gorm.DB, transactionId string) error {
	return tx.
		Where("id = ?", transactionId).
		Delete(models.Transaction{}).
		Error
}

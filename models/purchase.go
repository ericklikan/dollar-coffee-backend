package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model

	UserId     uuid.UUID      `gorm:"column:user_id;not null"`
	Items      []PurchaseItem `gorm:"foreignkey:transaction_id"`
	AmountPaid float32        `gorm:"type:decimal(12,2);not null"`
}

type PurchaseItem struct {
	gorm.Model

	TransactionId uint `gorm:"column:transaction_id;not null"`
	CoffeeItem    Coffee
	TypeOption    string `gorm:"type:varchar"`
}

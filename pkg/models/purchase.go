package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model

	UserId     uuid.UUID       `gorm:"column:user_id;not null"`
	Items      []*PurchaseItem `gorm:"foreignkey:transaction_id;PRELOAD:true"`
	AmountPaid float64         `gorm:"type:decimal(12,2);not null"`
	Total      float64         `gorm:"type:decimal(12,2);not null"`
}

type PurchaseItem struct {
	gorm.Model `json:"-"`

	TransactionId uint    `gorm:"column:transaction_id;not null" json:"-"`
	CoffeeId      uint    `gorm:"column:coffee_id;not null"`
	Price         float64 `gorm:"type:decimal(12,2);not null" json:"price"`
	TypeOption    string  `gorm:"type:text" json:"options"` //americano, latte, pourover, espresso, additional sugar/milk/cream
}

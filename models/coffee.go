package models

import "github.com/jinzhu/gorm"

type Coffee struct {
	gorm.Model

	Name        string  `json:"name" gorm:"type:varchar(255);not null"`
	Price       float32 `json:"price" gorm:"type:decimal(12,2);not null"`
	Description string  `json:"description" gorm:"type:text"`
}

type CoffeeType struct {
	gorm.Model

	Type string `gorm:"type:varchar(255);not null"`
}
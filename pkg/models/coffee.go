package models

import "github.com/jinzhu/gorm"

type Coffee struct {
	gorm.Model

	Name        string  `json:"name" gorm:"type:varchar(255);not null;unique_index"`
	Price       float64 `json:"price" gorm:"type:decimal(12,2);not null"`
	Description string  `json:"description" gorm:"type:text"`
	InStock     bool    `json:"inStock" gorm:"type:boolean;default:true"`
}

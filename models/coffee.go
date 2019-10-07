package models

import "github.com/jinzhu/gorm"

type Coffee struct {
	gorm.Model

	Name        string  `json:"name" gorm:"type:varchar(255);not null;unique_index"`
	Price       float32 `json:"price" gorm:"type:decimal(12,2);not null"`
	Description string  `json:"description" gorm:"type:text"`
}

func (coffee *Coffee) Create(db *gorm.DB) error {
	err := db.Create(coffee).Error
	if err != nil {
		return err
	}

	return nil
}

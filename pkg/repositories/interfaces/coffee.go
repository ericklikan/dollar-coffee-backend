package repository_interfaces

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

type CoffeePageQuery struct {
	PageQuery

	InStock *bool
}

type CoffeeRepository interface {
	CreateCoffee(tx *gorm.DB, coffee *models.Coffee) error
	GetCoffeesByIds(tx *gorm.DB, coffeeIds []string) (map[string]*models.Coffee, error)
	GetCoffeesPaginated(tx *gorm.DB, query *CoffeePageQuery) ([]*models.Coffee, error)
	UpdateCoffee(tx *gorm.DB, coffee *models.Coffee) error
	DeleteCoffee(tx *gorm.DB, coffeeId string) error
}

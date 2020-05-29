package repository

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/jinzhu/gorm"
)

type CoffeeRepositoryImpl struct {
	db *gorm.DB
}

func NewCoffeeRepository(db *gorm.DB) repository_interfaces.CoffeeRepository {
	return &CoffeeRepositoryImpl{
		db: db,
	}
}

func (repo *CoffeeRepositoryImpl) CreateCoffee(tx *gorm.DB, coffee *models.Coffee) error {
	return nil
}

func (repo *CoffeeRepositoryImpl) GetCoffeesByIds(tx *gorm.DB, coffeeIds []string) (map[string]*models.Coffee, error) {
	return nil, nil
}

func (repo *CoffeeRepositoryImpl) GetCoffeesPaginated(tx *gorm.DB, query *repository_interfaces.CoffeePageQuery) ([]*models.Coffee, error) {
	return nil, nil
}

func (repo *CoffeeRepositoryImpl) UpdateCoffee(tx *gorm.DB, coffee *models.Coffee) error {
	return nil
}

func (repo *CoffeeRepositoryImpl) DeleteCoffee(tx *gorm.DB, coffeeId string) error {
	return nil
}

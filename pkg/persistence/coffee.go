package persistence

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

func CreateCoffee(tx *gorm.DB, coffee *models.Coffee) error {
	return tx.Create(coffee).Error
}

func GetCoffeesByID(tx *gorm.DB, coffeeIds []string) (map[string]*models.Coffee, error) {
	var coffees []*models.Coffee
	if err := tx.
		Where("id in (?)", coffeeIds).
		Find(coffees).Error; err != nil {
		return nil, err
	}
	coffeesMap := make(map[string]*models.Coffee)
	for _, coffee := range coffees {
		coffeesMap[string(coffee.ID)] = coffee
	}
	return coffeesMap, nil
}

func GetCoffeesPaginated(tx *gorm.DB, pageSize int, page int, inStock *bool) ([]*models.Coffee, error) {
	var coffees []*models.Coffee
	q := tx.Model(models.Coffee{}).
		Select([]string{"ID", "name", "description", "price", "in_stock"}).
		Offset(page * pageSize).
		Limit(pageSize).
		Order("updated_at ASC")

	if inStock != nil {
		q = q.Where("in_stock = ?", *inStock)
	}
	if err := q.Find(&coffees).Error; err != nil {
		return nil, err
	}

	return coffees, nil
}

func UpdateCoffee(tx *gorm.DB, coffee *models.Coffee) error {
	return tx.Save(coffee).Error
}

func DeleteCoffee(tx *gorm.DB, coffeeId string) error {
	return tx.
		Where("id = ?", coffeeId).
		Delete(models.Coffee{}).
		Error
}

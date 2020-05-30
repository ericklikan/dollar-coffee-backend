package repository

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/ericklikan/dollar-coffee-backend/pkg/persistence"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const redisExpiryTime = 24
const redisMenuKey = "menu"

type CoffeeRepositoryImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCoffeeRepository(db *gorm.DB, redis *redis.Client) repository_interfaces.CoffeeRepository {
	return &CoffeeRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (repo *CoffeeRepositoryImpl) CreateCoffee(tx *gorm.DB, coffee *models.Coffee) error {
	// invalidate cache
	repo.redis.Del(redisMenuKey)
	return persistence.CreateCoffee(tx, coffee)
}

func (repo *CoffeeRepositoryImpl) GetCoffeesByIds(tx *gorm.DB, coffeeIds []string) (map[string]*models.Coffee, error) {
	return persistence.GetCoffeesByID(tx, coffeeIds)
}

func (repo *CoffeeRepositoryImpl) GetCoffeesPaginated(tx *gorm.DB, query *repository_interfaces.CoffeePageQuery) ([]*models.Coffee, error) {
	logger := log.WithFields(log.Fields{
		"Repository": "CoffeeRepository",
	})

	// retrieve redis key
	jsonString, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(jsonString))

	coffeePage, err := repo.redis.HGet(redisMenuKey, encodedQuery).Result()
	if err != nil || len(coffeePage) == 0 {
		// Don't fail request if redis doesn't work
		logger.WithError(err).Warn("Redis Error")

		newCoffeePage, err := persistence.GetCoffeesPaginated(tx, query.PageSize, query.Page, query.InStock)
		if err != nil {
			return nil, err
		}
		newCoffeePageJson, err := json.Marshal(newCoffeePage)
		if err != nil {
			logger.WithError(err).Warn("Error marshalling struct to json")
			return newCoffeePage, nil
		}

		repo.redis.HSet(redisMenuKey, encodedQuery, newCoffeePageJson)
		repo.redis.Expire(redisMenuKey, redisExpiryTime*time.Hour)
		return newCoffeePage, nil
	}

	logger.Debug("Fetching from redis")
	var coffeePageResult []*models.Coffee
	err = json.Unmarshal([]byte(coffeePage), &coffeePageResult)
	if err != nil {
		logger.WithError(err).Warn("Error unmarshalling page from json")
		return nil, err
	}

	return coffeePageResult, nil
}

func (repo *CoffeeRepositoryImpl) UpdateCoffee(tx *gorm.DB, coffee *models.Coffee) error {
	// invalidate cache
	repo.redis.Del(redisMenuKey)
	return persistence.UpdateCoffee(tx, coffee)
}

func (repo *CoffeeRepositoryImpl) DeleteCoffee(tx *gorm.DB, coffeeId string) error {
	// invalidate cache
	repo.redis.Del(redisMenuKey)
	return nil
}

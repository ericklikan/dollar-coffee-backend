package repository

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/ericklikan/dollar-coffee-backend/pkg/persistence"
	repository_interfaces "github.com/ericklikan/dollar-coffee-backend/pkg/repositories/interfaces"
	"github.com/jinzhu/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository_interfaces.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (repo *UserRepositoryImpl) GetUserByEmail(tx *gorm.DB, email string) (*models.User, error) {
	return persistence.GetUserByEmail(tx, email)
}

func (repo *UserRepositoryImpl) CreateUser(tx *gorm.DB, user *models.User) error {
	return persistence.CreateUser(tx, user)
}

func (repo *UserRepositoryImpl) GetUsersByIds(tx *gorm.DB, userIds []string) (map[string]*models.User, error) {
	return persistence.GetUsersByID(tx, userIds)
}

func (repo *UserRepositoryImpl) GetUsersPaginated(tx *gorm.DB, query *repository_interfaces.UsersPageQuery) ([]*models.User, error) {
	return persistence.GetUsersPaginated(tx, query.PageSize, query.Page)
}

func (repo *UserRepositoryImpl) UpdateUser(tx *gorm.DB, user *models.User) error {
	return persistence.UpdateUser(tx, user)
}

func (repo *UserRepositoryImpl) DeleteUser(tx *gorm.DB, userId string) error {
	return persistence.DeleteUser(tx, userId)
}

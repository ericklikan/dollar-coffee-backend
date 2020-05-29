package repository_interfaces

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

type UsersPageQuery struct {
	PageQuery
}

type UserRepository interface {
	LoginUser(email string, password string) *models.User
	CreateUser(tx *gorm.DB, user *models.User) error
	GetUsersPaginated(tx *gorm.DB, query UsersPageQuery) ([]*models.User, error)
	UpdateUser(tx *gorm.DB, coffee *models.User) error
	DeleteUser(tx *gorm.DB, userId string) error
}

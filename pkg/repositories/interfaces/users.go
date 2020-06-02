package repository_interfaces

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

type UsersPageQuery struct {
	PageQuery

	Role *string
}

type UserRepository interface {
	CreateUser(tx *gorm.DB, user *models.User) error
	GetUserByEmail(tx *gorm.DB, email string) (*models.User, error)
	GetUsersByIds(tx *gorm.DB, userIds []string) (map[string]*models.User, error)
	GetUsersPaginated(tx *gorm.DB, query *UsersPageQuery) ([]*models.User, error)
	UpdateUser(tx *gorm.DB, coffee *models.User) error
	DeleteUser(tx *gorm.DB, userId string) error
}

package persistence

import (
	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/jinzhu/gorm"
)

func CreateUser(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func GetUserByEmail(tx *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := tx.
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsersByID(tx *gorm.DB, userIds []string) (map[string]*models.User, error) {
	var users []*models.User
	if err := tx.
		Where("id in (?)", userIds).
		Find(&users).Error; err != nil {
		return nil, err
	}
	usersMap := make(map[string]*models.User)
	for _, user := range users {
		usersMap[user.ID.String()] = user
	}
	return usersMap, nil
}

func GetUsersPaginated(tx *gorm.DB, pageSize int, page int, role *string) ([]*models.User, error) {
	var users []*models.User
	q := tx.
		Offset(page * pageSize).
		Limit(pageSize)

	if role != nil {
		q = q.Where("role = ?", *role)
	}

	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

func UpdateUser(tx *gorm.DB, user *models.User) error {
	return tx.Save(user).Error
}

func DeleteUser(tx *gorm.DB, userId string) error {
	return tx.
		Where("id = ?", userId).
		Delete(models.User{}).
		Error
}

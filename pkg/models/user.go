package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"gopkg.in/go-playground/validator.v9"
)

type Token struct {
	jwt.StandardClaims

	UserId uuid.UUID

	// can be either user or admin
	Role string
}

type User struct {
	ID        uuid.UUID `gorm:"primary_key;column:id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FirstName   string  `json:"firstName" gorm:"type:varchar(255);not null"`
	LastName    string  `json:"lastName" gorm:"type:varchar(255);not null"`
	Email       string  `json:"email" gorm:"type:varchar(320);not null;unique_index" validate:"required,email"`
	PhoneNumber *string `json:"phoneNumber" gorm:"type:char(9)"`

	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty" gorm:"-"`

	// Role
	Role string `gorm:"type:varchar(10);default:'user'"`
}

func (user *User) Validate() error {
	v := validator.New()
	err := v.Struct(user)
	if err != nil {
		return errors.New("Invalid Email Address")
	}

	if len(user.Password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}

	return nil
}

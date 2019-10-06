package models

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

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

	Password string `json:"password"`
	Token    string `json:"token" gorm:"-"`

	// Role
	Role string `gorm:"type:varchar(10);default:'user'"`
}

func (user *User) Validate(db *gorm.DB) error {
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

func (user *User) Create(db *gorm.DB) (*User, error) {
	if err := user.Validate(db); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	dbTxInfo := db.Create(user)
	if dbTxInfo.Error != nil {
		return nil, err
	}

	//Create new JWT token for the newly registered account and default to role type as user
	tk := &Token{
		UserId: user.ID,
		Role:   "user",
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	return user, nil
}

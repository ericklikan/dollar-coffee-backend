package persistence

import (
	"testing"

	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/ericklikan/dollar-coffee-backend/pkg/test"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testUser := models.User{
		ID:        testUserId,
		FirstName: "Test",
		LastName:  "test",
		Email:     "test@testtest.test",
	}

	err := CreateUser(tx, &testUser)
	require.NoError(t, err)

	var retrievedUser models.User
	err = tx.Where("id = ?", testUser.ID).First(&retrievedUser).Error
	require.NoError(t, err)

	assert.Equal(t, testUser.FirstName, retrievedUser.FirstName)
	assert.Equal(t, testUser.Email, retrievedUser.Email)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestGetUsersById(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testUser := models.User{
		ID:        testUserId,
		FirstName: "Test",
		LastName:  "test",
		Email:     "test@testtest.test",
	}

	err := CreateUser(tx, &testUser)
	require.NoError(t, err)

	var retrievedUser models.User
	err = tx.Where("id = ?", testUser.ID).First(&retrievedUser).Error
	require.NoError(t, err)

	UsersMap, err := GetUsersByID(tx, []string{retrievedUser.ID.String()})
	require.NoError(t, err)

	User, doesUserExist := UsersMap[retrievedUser.ID.String()]
	assert.True(t, doesUserExist)
	assert.Equal(t, testUser.FirstName, User.FirstName)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestGetUsersPaginated(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testUser := models.User{
		ID:        testUserId,
		FirstName: "Test",
		LastName:  "test",
		Email:     "test@testtest.test",
	}

	err := CreateUser(tx, &testUser)
	require.NoError(t, err)

	page, err := GetUsersPaginated(tx, 1, 0, nil)
	require.NoError(t, err)
	assert.Len(t, page, 1)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestUpdateUser(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testUser := models.User{
		ID:        testUserId,
		FirstName: "Test",
		LastName:  "test",
		Email:     "test@testtest.test",
	}

	err := CreateUser(tx, &testUser)
	require.NoError(t, err)

	var retrievedUser models.User
	err = tx.Where("id = ?", testUser.ID).First(&retrievedUser).Error
	require.NoError(t, err)

	retrievedUser.FirstName = "Updated Name"
	err = UpdateUser(tx, &retrievedUser)
	require.NoError(t, err)

	err = tx.Where("id = ?", testUser.ID).First(&retrievedUser).Error
	require.NoError(t, err)

	assert.Equal(t, "Updated Name", retrievedUser.FirstName)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestDeleteUser(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testUser := models.User{
		ID:        testUserId,
		FirstName: "Test",
		LastName:  "test",
		Email:     "test@testtest.test",
	}

	err := CreateUser(tx, &testUser)
	require.NoError(t, err)

	var retrievedUser models.User
	err = tx.Where("id = ?", testUser.ID).First(&retrievedUser).Error
	require.NoError(t, err)

	err = DeleteUser(tx, testUser.ID.String())
	require.NoError(t, err)

	err = tx.Where("id = ?", testUser.ID).First(&retrievedUser).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

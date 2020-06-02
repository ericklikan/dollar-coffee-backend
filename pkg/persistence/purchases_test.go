package persistence

import (
	"strconv"
	"testing"

	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/ericklikan/dollar-coffee-backend/pkg/test"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testUserId, _ = uuid.Parse("594cb7e0-60a0-479d-ab2c-3a199d56cfee")

func TestCreateTransaction(t *testing.T) {
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

	testTransaction := models.Transaction{
		UserId:     testUserId,
		AmountPaid: 0,
		Total:      1.2,
	}
	testTransaction.ID = 735799

	err = CreateTransaction(tx, &testTransaction)
	require.NoError(t, err)

	var retrievedTransaction models.Transaction
	err = tx.Where("id = ?", testTransaction.ID).First(&retrievedTransaction).Error
	require.NoError(t, err)

	assert.Equal(t, testTransaction.UserId, retrievedTransaction.UserId)
	assert.Equal(t, testTransaction.Total, retrievedTransaction.Total)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestGetTransactionsById(t *testing.T) {
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

	testTransaction := models.Transaction{
		UserId:     testUserId,
		AmountPaid: 0,
		Total:      1.2,
	}
	testTransaction.ID = 735799

	err = CreateTransaction(tx, &testTransaction)
	require.NoError(t, err)

	var retrievedTransaction models.Transaction
	err = tx.Where("id = ?", testTransaction.ID).First(&retrievedTransaction).Error
	require.NoError(t, err)

	TransactionsMap, err := GetTransactionsByID(tx, []string{strconv.FormatUint(uint64(retrievedTransaction.ID), 10)})
	require.NoError(t, err)

	Transaction, doesTransactionExist := TransactionsMap[strconv.FormatUint(uint64(retrievedTransaction.ID), 10)]
	assert.True(t, doesTransactionExist)
	assert.Equal(t, testTransaction.UserId, Transaction.UserId)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestGetTransactionsPaginated(t *testing.T) {
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

	testTransaction := models.Transaction{
		UserId:     testUserId,
		AmountPaid: 0,
		Total:      1.2,
	}
	testTransaction.ID = 735799

	err = CreateTransaction(tx, &testTransaction)
	require.NoError(t, err)

	page, err := GetTransactionsPaginated(tx, 1, 0, nil, nil, nil)
	require.NoError(t, err)
	assert.Len(t, page, 1)

	thisShouldntExist := uuid.New().String()
	page, err = GetTransactionsPaginated(tx, 1, 0, &thisShouldntExist, nil, nil)
	require.NoError(t, err)
	assert.Len(t, page, 0)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestUpdateTransaction(t *testing.T) {
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

	testTransaction := models.Transaction{
		UserId:     testUserId,
		AmountPaid: 0,
		Total:      1.2,
	}
	testTransaction.ID = 735799

	err = CreateTransaction(tx, &testTransaction)
	require.NoError(t, err)

	var retrievedTransaction models.Transaction
	err = tx.Where("id = ?", testTransaction.ID).First(&retrievedTransaction).Error
	require.NoError(t, err)

	retrievedTransaction.AmountPaid = 1
	err = UpdateTransaction(tx, &retrievedTransaction)
	require.NoError(t, err)

	err = tx.Where("id = ?", testTransaction.ID).First(&retrievedTransaction).Error
	require.NoError(t, err)

	assert.Equal(t, float64(1), retrievedTransaction.AmountPaid)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestDeleteTransaction(t *testing.T) {
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

	testTransaction := models.Transaction{
		UserId:     testUserId,
		AmountPaid: 0,
		Total:      1.2,
	}
	testTransaction.ID = 735799

	err = CreateTransaction(tx, &testTransaction)
	require.NoError(t, err)

	var retrievedTransaction models.Transaction
	err = tx.Where("id = ?", testTransaction.ID).First(&retrievedTransaction).Error
	require.NoError(t, err)

	err = DeleteTransaction(tx, strconv.FormatUint(uint64(retrievedTransaction.ID), 10))
	require.NoError(t, err)

	err = tx.Where("id = ?", testTransaction.ID).First(&retrievedTransaction).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

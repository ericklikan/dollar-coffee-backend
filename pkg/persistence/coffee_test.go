package persistence

import (
	"strconv"
	"testing"

	"github.com/ericklikan/dollar-coffee-backend/pkg/models"
	"github.com/ericklikan/dollar-coffee-backend/pkg/test"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCoffee(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testCoffee := models.Coffee{
		Name:  "Test Coffee",
		Price: 1.2,
	}
	testCoffee.ID = 735799

	err := CreateCoffee(tx, &testCoffee)
	require.NoError(t, err)

	var retrievedCoffee models.Coffee
	err = tx.Where("id = ?", testCoffee.ID).First(&retrievedCoffee).Error
	require.NoError(t, err)

	assert.Equal(t, testCoffee.Name, retrievedCoffee.Name)
	assert.Equal(t, testCoffee.Price, retrievedCoffee.Price)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestGetCoffeesById(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testCoffee := models.Coffee{
		Name:  "Test Coffee",
		Price: 1.2,
	}
	testCoffee.ID = 735799

	err := CreateCoffee(tx, &testCoffee)
	require.NoError(t, err)

	var retrievedCoffee models.Coffee
	err = tx.Where("id = ?", testCoffee.ID).First(&retrievedCoffee).Error
	require.NoError(t, err)

	coffeesMap, err := GetCoffeesByID(tx, []string{strconv.FormatUint(uint64(retrievedCoffee.ID), 10)})
	require.NoError(t, err)

	coffee, doesCoffeeExist := coffeesMap[strconv.FormatUint(uint64(retrievedCoffee.ID), 10)]
	assert.True(t, doesCoffeeExist)
	assert.Equal(t, testCoffee.Name, coffee.Name)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestGetCoffeesPaginated(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testCoffee := models.Coffee{
		Name:  "Test Coffee",
		Price: 1.2,
	}
	testCoffee.ID = 735799

	err := CreateCoffee(tx, &testCoffee)
	require.NoError(t, err)

	testCoffee2 := models.Coffee{
		Name:  "Test Coffee 2",
		Price: 1.3,
	}
	testCoffee2.ID = 735800

	err = CreateCoffee(tx, &testCoffee2)
	require.NoError(t, err)

	page, err := GetCoffeesPaginated(tx, 1, 0, nil)
	require.NoError(t, err)
	assert.Len(t, page, 1)

	inStock := true
	page, err = GetCoffeesPaginated(tx, 1, 0, &inStock)
	require.NoError(t, err)
	require.Len(t, page, 1)
	assert.True(t, page[0].InStock)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestUpdateCoffee(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testCoffee := models.Coffee{
		Name:  "Test Coffee",
		Price: 1.2,
	}
	testCoffee.ID = 735799

	err := CreateCoffee(tx, &testCoffee)
	require.NoError(t, err)

	var retrievedCoffee models.Coffee
	err = tx.Where("id = ?", testCoffee.ID).First(&retrievedCoffee).Error
	require.NoError(t, err)

	retrievedCoffee.Name = "Updated Name"
	err = UpdateCoffee(tx, &retrievedCoffee)
	require.NoError(t, err)

	err = tx.Where("id = ?", testCoffee.ID).First(&retrievedCoffee).Error
	require.NoError(t, err)

	assert.Equal(t, "Updated Name", retrievedCoffee.Name)
	assert.Equal(t, testCoffee.Price, retrievedCoffee.Price)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

func TestDeleteCoffee(t *testing.T) {
	db := test.SetupTestDb(t)
	tx := db.Begin()

	testCoffee := models.Coffee{
		Name:  "Test Coffee",
		Price: 1.2,
	}
	testCoffee.ID = 735799

	err := CreateCoffee(tx, &testCoffee)
	require.NoError(t, err)

	var retrievedCoffee models.Coffee
	err = tx.Where("id = ?", testCoffee.ID).First(&retrievedCoffee).Error
	require.NoError(t, err)

	err = DeleteCoffee(tx, strconv.FormatUint(uint64(retrievedCoffee.ID), 10))
	require.NoError(t, err)

	err = tx.Where("id = ?", testCoffee.ID).First(&retrievedCoffee).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// rollback create because we don't want it to be in our db
	tx.Rollback()
}

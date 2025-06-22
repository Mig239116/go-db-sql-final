package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// setupTestDB производит подготвку тестовой базы
func setupTestDB(t *testing.T) (*sql.DB, ParcelStore) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	require.NoError(t, err)
	return db, NewParcelStore(db)
}

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, store := setupTestDB(t)
	defer db.Close()
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcelDb, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, id, parcelDb.Number)
	assert.Equal(t, parcel.Client, parcelDb.Client)
	assert.Equal(t, parcel.Status, parcelDb.Status)
	assert.Equal(t, parcel.Address, parcelDb.Address)
	assert.Equal(t, parcel.CreatedAt, parcelDb.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.Get(id)
	require.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, store := setupTestDB(t)
	defer db.Close()
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	newParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, newParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {

	db, store := setupTestDB(t)
	defer db.Close()
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	newParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, newParcel.Status)
}

func TestGetByClient(t *testing.T) {
	// prepare
	db, store := setupTestDB(t)
	defer db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)
		assert.Equal(t, i+1, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		expectedParcel, exists := parcelMap[parcel.Number]
		assert.True(t, exists)
		assert.Equal(t, expectedParcel, parcel)
	}
}

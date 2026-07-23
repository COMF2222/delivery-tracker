//go:build integration

package repository

import (
	"delivery-tracker/internal/domain"
	"delivery-tracker/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParcelRepository_CreateParcel(t *testing.T) {
	db := testhelpers.NewPostgresContainer(t)
	repo := NewParcelRepository(db)

	t.Run("success", func(t *testing.T) {
		parcel := domain.Parcel{
			TrackNumber:      "Q4P405SHH8EG",
			ItemName:         "Reina Petite Ring",
			RecipientName:    "Иван",
			RecipientPhone:   "79999999999",
			RecipientAddress: "Moscow",
			CurrentStatus:    domain.StatusCreated,
			CurrentLocation:  "Berlin",
			IsArchived:       false,
		}

		err := repo.CreateParcel(&parcel, 1)
		require.NoError(t, err)

		saved, err := repo.GetByID(parcel.ID)
		require.NoError(t, err)
		assert.Greater(t, saved.ID, 0)
		assert.Equal(t, parcel.ID, saved.ID)
		assert.Equal(t, "Q4P405SHH8EG", saved.TrackNumber)
		assert.Equal(t, "Reina Petite Ring", saved.ItemName)
		assert.Equal(t, "Иван", saved.RecipientName)
		assert.Equal(t, "79999999999", saved.RecipientPhone)
		assert.Equal(t, "Moscow", saved.RecipientAddress)
		assert.Equal(t, domain.StatusCreated, saved.CurrentStatus)
		assert.Equal(t, "Berlin", saved.CurrentLocation)
		assert.False(t, saved.IsArchived)
	})

	t.Run("track number conflict", func(t *testing.T) {
		parcel := domain.Parcel{
			TrackNumber:      "Q4P405SHH7EG",
			ItemName:         "Reina Petite Ring",
			RecipientName:    "Иван",
			RecipientPhone:   "79999999999",
			RecipientAddress: "Moscow",
			CurrentStatus:    domain.StatusCreated,
			CurrentLocation:  "Berlin",
			IsArchived:       false,
		}
		err := repo.CreateParcel(&parcel, 1)
		require.NoError(t, err)

		parcelDuplicate := domain.Parcel{
			TrackNumber:      "Q4P405SHH7EG",
			ItemName:         "Reina Petite Ring",
			RecipientName:    "Иван",
			RecipientPhone:   "79999999999",
			RecipientAddress: "Moscow",
			CurrentStatus:    domain.StatusCreated,
			CurrentLocation:  "Berlin",
			IsArchived:       false,
		}
		err = repo.CreateParcel(&parcelDuplicate, 1)
		assert.ErrorIs(t, err, ErrTrackNumberAlreadyExists)
	})

	t.Run("create error", func(t *testing.T) {
		parcel := domain.Parcel{
			TrackNumber:      "Q4P403SHH7EG",
			ItemName:         "Reina Petite Ring",
			RecipientName:    "Иван",
			RecipientPhone:   "79999999999",
			RecipientAddress: "Moscow",
			CurrentStatus:    domain.StatusCreated,
			CurrentLocation:  "Berlin",
			IsArchived:       false,
		}

		err := repo.CreateParcel(&parcel, 9999999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create parcel")
	})
}

func TestParcelRepository_GetByTrackNumber(t *testing.T) {
	db := testhelpers.NewPostgresContainer(t)
	repo := NewParcelRepository(db)

	trackNumber := "Q4P405SHH8EG"

	err := repo.CreateParcel(&domain.Parcel{
		TrackNumber:      "Q4P405SHH8EG",
		ItemName:         "Reina Petite Ring",
		RecipientName:    "Иван",
		RecipientPhone:   "79999999999",
		RecipientAddress: "Moscow",
		CurrentStatus:    domain.StatusCreated,
		CurrentLocation:  "Berlin",
		IsArchived:       false,
	}, 1)
	require.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		parcel, err := repo.GetByTrackNumber(trackNumber)
		require.NoError(t, err)
		assert.Equal(t, trackNumber, parcel.TrackNumber)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByTrackNumber("")
		assert.ErrorIs(t, err, ErrParcelNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		err := db.Close()
		require.NoError(t, err)

		_, err = repo.GetByTrackNumber(trackNumber)
		require.Error(t, err)
		assert.NotErrorIs(t, err, ErrParcelNotFound)
		assert.Contains(t, err.Error(), "failed to get parcel by track number")
	})
}

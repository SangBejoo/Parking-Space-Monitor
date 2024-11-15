package repository

import (
	"database/sql"
)

// Repository struct to hold all repositories

type Repository struct {
	DB *sql.DB

	TaxiRepository *TaxiRepository

	PlaceRepository *PlaceRepository

	MappingRepository *MappingRepository

	CountersRepository *MappingRepository
}

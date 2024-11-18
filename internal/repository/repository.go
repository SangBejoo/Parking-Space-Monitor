package repository

import (
	"database/sql"
	"/internal/repository"
)

// Repository struct to hold all repositories

type Repository struct {
	DB *sql.DB

	TaxiRepository *TaxiRepository

	PlaceRepository *PlaceRepository

	MappingRepository *MappingRepository
	CountersRepository *repository.MappingRepository
	CountersRepository *MappingRepository
}

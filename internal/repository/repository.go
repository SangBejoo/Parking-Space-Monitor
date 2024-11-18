package repository

import (
    "database/sql"
)

// CountersRepository handles operations related to counters.
type CountersRepository struct {
    DB *sql.DB
}

type Repository struct {
    DB                  *sql.DB
    TaxiRepository      *TaxiRepository
    PlaceRepository     *PlaceRepository
    MappingRepository   *MappingRepository
    CountersRepository  *CountersRepository 
}
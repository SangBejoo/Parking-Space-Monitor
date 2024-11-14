// internal/repository/mapping_repository.go
package repository

import "database/sql"

// MappingRepository handles operations related to mappings and counters.
type MappingRepository struct {
    DB *sql.DB
}

// InsertMapping inserts a new mapping into the mapping table.
func (mr *MappingRepository) InsertMapping(taxiID string, placeID int) error {
    _, err := mr.DB.Exec("INSERT INTO mapping (taxi_id, place_id) VALUES ($1, $2)", taxiID, placeID)
    return err
}

// GetCounter retrieves the current counter for a taxi and place.
func (mr *MappingRepository) GetCounter(taxiID string, placeID int) (int, error) {
    var count int
    err := mr.DB.QueryRow("SELECT counter FROM counters WHERE taxi_id = $1 AND place_id = $2", taxiID, placeID).Scan(&count)
    return count, err
}

// InsertCounter inserts a new counter.
func (mr *MappingRepository) InsertCounter(taxiID string, placeID int) error {
    _, err := mr.DB.Exec("INSERT INTO counters (taxi_id, place_id, counter, last_counted) VALUES ($1, $2, 1, CURRENT_TIMESTAMP)", taxiID, placeID)
    return err
}

// UpdateCounter increments the counter.
func (mr *MappingRepository) UpdateCounter(taxiID string, placeID int) error {
    _, err := mr.DB.Exec("UPDATE counters SET counter = counter + 1, last_counted = CURRENT_TIMESTAMP WHERE taxi_id = $1 AND place_id = $2", taxiID, placeID)
    return err
}
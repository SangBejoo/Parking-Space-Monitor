// internal/repository/mapping_repository.go
package repository

<<<<<<< HEAD
import (
	"database/sql"
	"fmt"
	"log"
)
=======
import "database/sql"
>>>>>>> parent of ec2908f (update)

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
<<<<<<< HEAD
	_, err := mr.DB.Exec("UPDATE counters SET counter = counter + 1, last_counted = CURRENT_TIMESTAMP WHERE taxi_id = $1 AND place_id = $2", taxiID, placeID)
	return err
}

// internal/repository/mapping_repository.go

// GetAllMappings retrieves all records from mapping table
func (mr *MappingRepository) GetAllMappings() ([]map[string]interface{}, error) {
	query := `
        SELECT m.taxi_id, m.place_id, p.place_name
        FROM mapping m
        JOIN places p ON m.place_id = p.place_id
    `

	rows, err := mr.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query mappings: %v", err)
	}
	defer rows.Close()

	var mappings []map[string]interface{}
	for rows.Next() {
		var taxiID string
		var placeID int
		var placeName string

		if err := rows.Scan(&taxiID, &placeID, &placeName); err != nil {
			return nil, fmt.Errorf("failed to scan mapping: %v", err)
		}

		mappings = append(mappings, map[string]interface{}{
			"taxi_id":    taxiID,
			"place_id":   placeID,
			"place_name": placeName,
		})
	}

	return mappings, nil
}

// UpsertMapping inserts a new mapping or updates it if it already exists.
func (mr *MappingRepository) UpsertMapping(taxiID string, placeID int) error {
    query := `
        INSERT INTO mapping (taxi_id, place_id)
        VALUES ($1, $2)
        ON CONFLICT (taxi_id, place_id) DO NOTHING
    `
    _, err := mr.DB.Exec(query, taxiID, placeID)
    if err != nil {
        log.Printf("Error upserting mapping for taxi %s and place %d: %v", taxiID, placeID, err)
    }
    return err
}

// UpdateTaxiDuration updates the duration of a taxi in a specific place.

func (r *MappingRepository) UpdateTaxiDuration(taxiID, placeID string) error {

	query := "UPDATE taxi_durations SET duration = duration + 1 WHERE taxi_id = ? AND place_id = ?"

	_, err := r.DB.Exec(query, taxiID, placeID)

	if err != nil {

		log.Printf("Error updating taxi duration: %v", err)

		return err

	}

	return nil

}

// ResetTaxiDuration resets the duration of a taxi in a place.

func (repo *MappingRepository) ResetTaxiDuration(taxiID string) error {
	query := "UPDATE taxi_mapping SET duration = 0 WHERE taxi_id = ?"

	_, err := repo.DB.Exec(query, taxiID)

	if err != nil {
		log.Printf("Error resetting taxi duration for taxi %s: %v", taxiID, err)
		return err
	}

	return nil
}

=======
    _, err := mr.DB.Exec("UPDATE counters SET counter = counter + 1, last_counted = CURRENT_TIMESTAMP WHERE taxi_id = $1 AND place_id = $2", taxiID, placeID)
    return err
}
>>>>>>> parent of ec2908f (update)

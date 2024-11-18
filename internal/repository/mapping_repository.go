// internal/repository/mapping_repository.go

package repository

import (
    "database/sql"
    "fmt"
   
    "github.com/SangBejoo/parking-space-monitor/internal/models"
)

// MappingRepository handles operations related to mappings.
type MappingRepository struct {
    DB *sql.DB
}

// InsertMapping inserts a new mapping into the mapping table.
func (mr *MappingRepository) InsertMapping(mapping models.Mapping) error {
    query := `
        INSERT INTO mapping (place_id, taxi_id)
        VALUES ($1, $2)
    `
    _, err := mr.DB.Exec(query, mapping.PlaceID, mapping.TaxiID)
    if err != nil {
        return fmt.Errorf("failed to insert mapping: %w", err)
    }
    return nil
}

// GetAllMappings retrieves all mappings with place names and taxi numbers.
func (mr *MappingRepository) GetAllMappings() ([]models.Mapping, error) {
    query := `
        SELECT m.id, p.place_name, t.nomor_taxi
        FROM mapping m
        JOIN places p ON m.place_id = p.place_id
        JOIN taxi t ON m.taxi_id = t.id
    `

    rows, err := mr.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query mappings: %v", err)
    }
    defer rows.Close()

    var mappings []models.Mapping
    for rows.Next() {
        var mapping models.Mapping
        if err := rows.Scan(&mapping.ID, &mapping.PlaceName, &mapping.TaxiID); err != nil {
            return nil, fmt.Errorf("failed to scan mapping: %v", err)
        }
        mappings = append(mappings, mapping)
    }

    return mappings, nil
}

// GetMappingByID retrieves a mapping by its ID.
func (mr *MappingRepository) GetMappingByID(mappingID int) (*models.Mapping, error) {
    var mapping models.Mapping
    query := `
        SELECT m.id, p.place_name, t.nomor_taxi
        FROM mapping m
        JOIN places p ON m.place_id = p.place_id
        JOIN taxi t ON m.taxi_id = t.id
        WHERE m.id = $1
    `
    err := mr.DB.QueryRow(query, mappingID).Scan(&mapping.ID, &mapping.PlaceName, &mapping.TaxiID)
    if err != nil {
        return nil, err
    }
    return &mapping, nil
}

// UpdateMapping updates an existing mapping.
func (mr *MappingRepository) UpdateMapping(mapping models.Mapping) error {
    query := `
        UPDATE mapping
        SET place_id = $1, taxi_id = $2
        WHERE id = $3
    `
    res, err := mr.DB.Exec(query, mapping.PlaceID, mapping.TaxiID, mapping.ID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("mapping not found")
    }

    return nil
}

// DeleteMapping deletes a mapping by its ID.
func (mr *MappingRepository) DeleteMapping(mappingID int) error {
    res, err := mr.DB.Exec("DELETE FROM mapping WHERE id = $1", mappingID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("mapping not found")
    }

    return nil
}

// UpdateTaxiDuration updates the duration a taxi has spent in a place.
func (mr *MappingRepository) UpdateTaxiDuration(taxiID string, placeID int) error {
    query := `UPDATE taxi_durations SET place_id = $1, updated_at = NOW() WHERE taxi_id = $2`
    _, err := mr.DB.Exec(query, placeID, taxiID)
    return err
}

// ResetTaxiDuration resets the duration a taxi has spent in any place.
func (mr *MappingRepository) ResetTaxiDuration(taxiID string) error {
    query := `UPDATE taxi_durations SET place_id = NULL, updated_at = NOW() WHERE taxi_id = $1`
    _, err := mr.DB.Exec(query, taxiID)
    return err
}
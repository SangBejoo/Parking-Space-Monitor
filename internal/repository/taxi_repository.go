// internal/repository/taxi_repository.go
package repository

import (
    "database/sql"
    "fmt"

    "github.com/SangBejoo/parking-space-monitor/internal/models"
)

// TaxiRepository handles CRUD operations for TaxiLocation.
type TaxiRepository struct {
    DB *sql.DB
}

// CreateTaxi inserts a new taxi location into the database.
func (tr *TaxiRepository) CreateTaxi(location models.TaxiLocation) error {
    _, err := tr.DB.Exec(`INSERT INTO taxi_location (taxi_id, longitude, latitude, updated_at) 
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP) 
        ON CONFLICT (taxi_id) DO NOTHING`,
        location.TaxiID, location.Longitude, location.Latitude)
    return err
}

// GetAllTaxis retrieves all taxi locations from the database.
func (tr *TaxiRepository) GetAllTaxis() ([]models.TaxiLocation, error) {
    rows, err := tr.DB.Query("SELECT taxi_id, longitude, latitude FROM taxi_location")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var taxis []models.TaxiLocation
    for rows.Next() {
        var taxi models.TaxiLocation
        if err := rows.Scan(&taxi.TaxiID, &taxi.Longitude, &taxi.Latitude); err != nil {
            return nil, err
        }
        taxis = append(taxis, taxi)
    }

    return taxis, nil
}

// GetTaxiByID retrieves a taxi location by its ID.
func (tr *TaxiRepository) GetTaxiByID(taxiID string) (*models.TaxiLocation, error) {
    var taxi models.TaxiLocation
    err := tr.DB.QueryRow("SELECT taxi_id, longitude, latitude FROM taxi_location WHERE taxi_id = $1", taxiID).
        Scan(&taxi.TaxiID, &taxi.Longitude, &taxi.Latitude)
    if err != nil {
        return nil, err
    }
    return &taxi, nil
}

// UpdateTaxi updates an existing taxi location.
func (tr *TaxiRepository) UpdateTaxi(taxiID string, location models.TaxiLocation) error {
    res, err := tr.DB.Exec(`UPDATE taxi_location SET longitude = $1, latitude = $2, updated_at = CURRENT_TIMESTAMP 
        WHERE taxi_id = $3`,
        location.Longitude, location.Latitude, taxiID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("taxi not found")
    }

    return nil
}

// DeleteTaxi deletes a taxi location by its ID.
func (tr *TaxiRepository) DeleteTaxi(taxiID string) error {
    res, err := tr.DB.Exec("DELETE FROM taxi_location WHERE taxi_id = $1", taxiID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("taxi not found")
    }

    return nil
}
// internal/repository/taxi_repository.go
package repository

import (
<<<<<<< HEAD
	"database/sql"
	"fmt"
	"log"
=======
    "database/sql"
    "fmt"
>>>>>>> parent of ec2908f (update)

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

// CreateMultipleTaxis inserts multiple taxi locations into the database.
func (tr *TaxiRepository) CreateMultipleTaxis(locations []models.TaxiLocation) error {
	tx, err := tr.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO taxi_location (taxi_id, longitude, latitude, updated_at) 
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP) 
        ON CONFLICT (taxi_id) DO UPDATE SET longitude = EXCLUDED.longitude, latitude = EXCLUDED.latitude, updated_at = CURRENT_TIMESTAMP`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, location := range locations {
		_, err := stmt.Exec(location.TaxiID, location.Longitude, location.Latitude)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
<<<<<<< HEAD
	log.Printf("Attempting to delete taxi with ID: %s", taxiID)

	res, err := tr.DB.Exec("DELETE FROM taxi_location WHERE taxi_id = $1", taxiID)
	if err != nil {
		log.Printf("Database error when deleting taxi %s: %v", taxiID, err)
		return fmt.Errorf("database error: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected for taxi %s: %v", taxiID, err)
		return fmt.Errorf("error checking deletion result: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("No taxi found with ID: %s", taxiID)
		return fmt.Errorf("taxi not found")
	}

	log.Printf("Successfully deleted taxi %s, rows affected: %d", taxiID, rowsAffected)
	return nil
}
=======
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
>>>>>>> parent of ec2908f (update)

// internal/repository/place_repository.go
package repository

import (
    "database/sql"
    "fmt"

    "github.com/SangBejoo/parking-space-monitor/internal/models"
)

// PlaceRepository handles CRUD operations for Place.
type PlaceRepository struct {
    DB *sql.DB
}

// CreatePlace inserts a new place into the database.
func (pr *PlaceRepository) CreatePlace(place models.Place) (int, error) {
    var placeID int
    err := pr.DB.QueryRow(`INSERT INTO places (place_name, polygon) 
        VALUES ($1, $2) RETURNING place_id`,
        place.PlaceName, place.Polygon).Scan(&placeID)
    if err != nil {
        return 0, err
    }
    return placeID, nil
}

// GetAllPlaces retrieves all places from the database.
func (pr *PlaceRepository) GetAllPlaces() ([]models.Place, error) {
    rows, err := pr.DB.Query("SELECT place_id, place_name, polygon FROM places")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var places []models.Place
    for rows.Next() {
        var place models.Place
        if err := rows.Scan(&place.PlaceID, &place.PlaceName, &place.Polygon); err != nil {
            return nil, err
        }
        places = append(places, place)
    }

    return places, nil
}

// GetPlaceByID retrieves a place by its ID.
func (pr *PlaceRepository) GetPlaceByID(placeID int) (*models.Place, error) {
    var place models.Place
    err := pr.DB.QueryRow("SELECT place_id, place_name, polygon FROM places WHERE place_id = $1", placeID).
        Scan(&place.PlaceID, &place.PlaceName, &place.Polygon)
    if err != nil {
        return nil, err
    }
    return &place, nil
}

// UpdatePlace updates an existing place.
func (pr *PlaceRepository) UpdatePlace(placeID int, place models.Place) error {
    res, err := pr.DB.Exec(`UPDATE places SET place_name = $1, polygon = $2 WHERE place_id = $3`,
        place.PlaceName, place.Polygon, placeID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("place not found")
    }

    return nil
}

// DeletePlace deletes a place by its ID.
func (pr *PlaceRepository) DeletePlace(placeID int) error {
    res, err := pr.DB.Exec("DELETE FROM places WHERE place_id = $1", placeID)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("place not found")
    }

    return nil
}
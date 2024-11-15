// internal/repository/place_repository.go
package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SangBejoo/parking-space-monitor/internal/models"
)

// PlaceRepository handles CRUD operations for Place.
type PlaceRepository struct {
	DB *sql.DB
}

// internal/repository/place_repository.go

// internal/repository/place_repository.go

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

func (pr *PlaceRepository) CreateMultiplePlaces(places []models.Place) error {
    tx, err := pr.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`INSERT INTO places (place_name, polygon) VALUES ($1, $2) RETURNING place_id`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for i, place := range places {
        err := stmt.QueryRow(place.PlaceName, place.Polygon).Scan(&places[i].PlaceID)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (pr *PlaceRepository) GetAllPlaces() ([]models.Place, error) {
	query := `
        SELECT 
            place_id, 
            place_name, 
            CASE 
                WHEN polygon IS NULL THEN '{"type":"Polygon","coordinates":[]}'::jsonb
                WHEN jsonb_typeof(polygon) = 'array' THEN 
                    jsonb_build_object(
                        'type', 'Polygon',
                        'coordinates', jsonb_build_array(polygon)
                    )
                ELSE polygon
            END as polygon
        FROM places
    `

	rows, err := pr.DB.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var places []models.Place
	for rows.Next() {
		var place models.Place
		var polygonBytes []byte

		if err := rows.Scan(&place.PlaceID, &place.PlaceName, &polygonBytes); err != nil {
			log.Printf("Row scan error: %v", err)
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		log.Printf("Raw polygon data: %s", string(polygonBytes))

		if err := json.Unmarshal(polygonBytes, &place.Polygon); err != nil {
			// Try to convert array format to GeoJSON
			var coordinates [][][]json.Number
			if err := json.Unmarshal(polygonBytes, &coordinates); err == nil {
				place.Polygon = models.GeoJSONPolygon{
					Type:        "Polygon",
					Coordinates: coordinates,
				}
			} else {
				log.Printf("Polygon unmarshal error: %v, data: %s", err, string(polygonBytes))
				return nil, fmt.Errorf("polygon unmarshal error: %v", err)
			}
		}

		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %v", err)
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

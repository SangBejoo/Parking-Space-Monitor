// internal/models/place.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type GeoJSONPolygon struct {
	Type        string            `json:"type"`
	Coordinates [][][]json.Number `json:"coordinates"`
}

// Scan implements sql.Scanner interface
func (p *GeoJSONPolygon) Scan(value interface{}) error {
	if value == nil {
		return errors.New("scanning nil value")
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	var geoJSON struct {
		Type        string            `json:"type"`
		Coordinates [][][]json.Number `json:"coordinates"`
	}

	if err := json.Unmarshal(b, &geoJSON); err == nil {
		p.Type = geoJSON.Type
		p.Coordinates = geoJSON.Coordinates
		return nil
	}

	// Try to unmarshal as array if GeoJSON fails
	var coords [][][]json.Number
	if err := json.Unmarshal(b, &coords); err == nil {
		p.Type = "Polygon"
		p.Coordinates = coords
		return nil
	}

	return fmt.Errorf("failed to unmarshal polygon data")
}

// Value implements driver.Valuer interface
func (p GeoJSONPolygon) Value() (driver.Value, error) {
	return json.Marshal(struct {
		Type        string            `json:"type"`
		Coordinates [][][]json.Number `json:"coordinates"`
	}{
		Type:        "Polygon",
		Coordinates: p.Coordinates,
	})
}

type Place struct {
	PlaceID   int            `json:"place_id"`
	PlaceName string         `json:"place_name"`
	Polygon   GeoJSONPolygon `json:"polygon"`
}

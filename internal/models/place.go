// internal/models/place.go
package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "fmt"
)

type GeoJSONPolygon struct {
    Type        string        `json:"type"`
    Coordinates [][][]float64 `json:"coordinates"`
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

    // First try direct array format
    var coords [][][]float64
    if err := json.Unmarshal(b, &coords); err == nil {
        p.Type = "Polygon"
        p.Coordinates = coords
        return nil
    }

    // Try full GeoJSON format
    type fullGeoJSON struct {
        Type        string        `json:"type"`
        Coordinates [][][]float64 `json:"coordinates"`
    }
    var gj fullGeoJSON
    if err := json.Unmarshal(b, &gj); err != nil {
        return fmt.Errorf("failed to unmarshal polygon data: %v", err)
    }

    p.Type = gj.Type
    p.Coordinates = gj.Coordinates
    return nil
}

// Value implements driver.Valuer interface
func (p GeoJSONPolygon) Value() (driver.Value, error) {
    return json.Marshal(struct {
        Type        string        `json:"type"`
        Coordinates [][][]float64 `json:"coordinates"`
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
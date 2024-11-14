// internal/models/place.go
package models

import "encoding/json"

// Place represents a geographical place with a polygon.
type Place struct {
    PlaceID   int             `json:"place_id"`
    PlaceName string          `json:"place_name"`
    Polygon   json.RawMessage `json:"polygon"` // GeoJSON Geometry
}

// GeoJSONGeometry represents the geometry part of a GeoJSON object.
type GeoJSONGeometry struct {
    Type        string        `json:"type"`
    Coordinates [][][]float64 `json:"coordinates"`
}

// internal/models/mapping.go

package models

// Mapping represents the association between a taxi and a place.
type Mapping struct {
    ID        int    `json:"id"`
    PlaceID   int    `json:"place_id"`
    TaxiID    int    `json:"taxi_id"`
    PlaceName string `json:"place_name,omitempty"` // For response purposes
}
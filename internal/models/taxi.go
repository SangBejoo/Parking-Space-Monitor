// internal/models/taxi.go
package models

// TaxiLocation represents the taxi's geographic location.
type TaxiLocation struct {
	TaxiID    string  `json:"taxi_id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

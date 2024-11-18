// internal/models/taxi.go
package models

// Taxi represents a taxi entity.
type Taxi struct {
    ID        int     `json:"id"`
    NomorTaxi string  `json:"nomor_taxi"`
    Longitude float64 `json:"longitude"`
    Latitude  float64 `json:"latitude"`
}

// TaxiLocation represents a taxi's location.
type TaxiLocation struct {
    TaxiID    string  `json:"taxi_id"`
    Longitude float64 `json:"longitude"`
    Latitude  float64 `json:"latitude"`
}
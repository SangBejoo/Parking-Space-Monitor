package scheduler

import (
    "log"
    "sync"

    "github.com/SangBejoo/parking-space-monitor/internal/repository"
)

// Scheduler handles scheduled tasks like mapping taxis to places.
type Scheduler struct {
    Repo  *repository.Repository
    Mutex sync.Mutex
}

// ProcessTaxi processes a taxi's location and updates it in the database.
func (s *Scheduler) ProcessTaxi(taxiID string, longitude, latitude float64) {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()

    log.Printf("Processing taxi %s at coordinates (%.6f, %.6f)", taxiID, longitude, latitude)

    // Update the taxi location in the database
    err := s.Repo.TaxiRepository.UpdateTaxiLocation(taxiID, longitude, latitude)
    if err != nil {
        log.Printf("Error updating taxi location: %v", err)
    }
}
// Point represents a geographic coordinate.
type Point struct {
    X float64 // longitude
    Y float64 // latitude
}

// isPointInPolygon checks if a point lies within a polygon using the ray-casting algorithm.
func isPointInPolygon(longitude, latitude float64, polygon []Point) bool {
    intersects := false
    j := len(polygon) - 1
    for i := 0; i < len(polygon); i++ {
        xi, yi := polygon[i].X, polygon[i].Y
        xj, yj := polygon[j].X, polygon[j].Y

        intersect := ((yi > latitude) != (yj > latitude)) &&
            (longitude < (xj-xi)*(latitude-yi)/(yj-yi)+xi)
        if intersect {
            intersects = !intersects
        }
        j = i
    }
    return intersects
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(repo *repository.Repository) *Scheduler {
    return &Scheduler{
        Repo: repo,
    }
}

// MapTaxiLocations assigns taxis to places based on their current locations.
func (s *Scheduler) MapTaxiLocations() {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()
    log.Println("Starting taxi location mapping...")

    taxis, err := s.Repo.TaxiRepository.GetAllTaxis()
    if err != nil {
        log.Printf("Error getting taxis: %v", err)
        return
    }

    places, err := s.Repo.PlaceRepository.GetAllPlaces()
    if err != nil {
        log.Printf("Error getting places: %v", err)
        return
    }

    for _, taxi := range taxis {
        log.Printf("Processing taxi %s at coordinates (%.2f, %.2f)",
            taxi.TaxiID, taxi.Longitude, taxi.Latitude)

        matched := false
        for _, place := range places {
            log.Printf("Checking against place %s with polygon: %v",
                place.PlaceName, place.Polygon)

            // Convert GeoJSONPolygon to []Point
            var polygon []Point
            if len(place.Polygon.Coordinates) > 0 {
                for _, coord := range place.Polygon.Coordinates[0] {
                    if len(coord) >= 2 {
                        lon, err1 := coord[0].Float64()
                        lat, err2 := coord[1].Float64()
                        if err1 != nil || err2 != nil {
                            log.Printf("Error converting coordinate to float64 for place %s", place.PlaceName)
                            continue
                        }
                        polygon = append(polygon, Point{
                            X: lon,
                            Y: lat,
                        })
                    }
                }
            } else {
                log.Printf("No coordinates found for place %s", place.PlaceName)
                continue
            }

            // Check if taxi is within polygon
            if isPointInPolygon(taxi.Longitude, taxi.Latitude, polygon) {
                matched = true
                log.Printf("Taxi %s is within %s", taxi.TaxiID, place.PlaceName)

                // Update taxi duration with placeID as int
                err := s.Repo.MappingRepository.UpdateTaxiDuration(taxi.TaxiID, place.PlaceID)
                if err != nil {
                    log.Printf("Error updating taxi duration: %v", err)
                }

                break
            }
        }

        if !matched {
            log.Printf("Taxi %s does not match any place", taxi.TaxiID)
            // Reset duration if taxi moved out
            err := s.Repo.MappingRepository.ResetTaxiDuration(taxi.TaxiID)
            if err != nil {
                log.Printf("Error resetting taxi duration: %v", err)
            }
        }
    }

    log.Println("Mapping process completed")
}
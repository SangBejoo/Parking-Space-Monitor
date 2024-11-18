// internal/scheduler/scheduler.go
package scheduler

import (
    "database/sql"
    "log"
    "sync"
    "strconv"

    "github.com/SangBejoo/parking-space-monitor/internal/repository"
)

// Scheduler handles scheduled tasks like mapping taxis to places.
type Scheduler struct {
    Repo  *Repository
    Mutex sync.Mutex
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

// Repository aggregates all repository dependencies.
type Repository struct {
    DB          *sql.DB
    TaxiRepo    *repository.TaxiRepository
    PlaceRepo   *repository.PlaceRepository
    MappingRepo *repository.MappingRepository
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(repo *Repository) *Scheduler {
    return &Scheduler{
        Repo: repo,
    }
}

// MapTaxiLocations assigns taxis to places based on their current locations.
func (s *Scheduler) MapTaxiLocations() {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()
    log.Println("Starting taxi location mapping...")

    taxis, err := s.Repo.TaxiRepo.GetAllTaxis()
    if err != nil {
        log.Printf("Error getting taxis: %v", err)
        return
    }

    places, err := s.Repo.PlaceRepo.GetAllPlaces()
    if err != nil {
        log.Printf("Error getting places: %v", err)
        return
    }

    for _, taxi := range taxis {
        log.Printf("Processing taxi %s at coordinates (%f, %f)",
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
                        lonStr := coord[0]
                        latStr := coord[1]
                        lon, err1 := lonStr.Float64()
                        lat, err2 := latStr.Float64()
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

                // Update taxi duration
                err := s.Repo.MappingRepo.UpdateTaxiDuration(taxi.TaxiID, strconv.Itoa(place.PlaceID))
                if err != nil {
                    log.Printf("Error updating taxi duration: %v", err)
                }

                break
            }
        }

        if !matched {
            log.Printf("Taxi %s does not match any place", taxi.TaxiID)
            // Reset duration if taxi moved out
            err := s.Repo.MappingRepo.ResetTaxiDuration(taxi.TaxiID)
            if err != nil {
                log.Printf("Error resetting taxi duration: %v", err)
            }
        }
    }

    log.Println("Mapping process completed")
}
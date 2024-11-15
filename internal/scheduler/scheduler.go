// internal/scheduler/scheduler.go
package scheduler

import (
    "database/sql"
    "sync"
    "log"
    "encoding/json"

    "github.com/SangBejoo/parking-space-monitor/internal/repository"
)

// Scheduler handles scheduled tasks like mapping taxis to places.
type Scheduler struct {
    Repo  *Repository
    Mutex sync.Mutex
}

// isPointInPolygon checks if a point is inside a polygon.
// In scheduler.go

// Point represents a geographic coordinate
type Point struct {
    X float64 // longitude
    Y float64 // latitude
}

// isPointInPolygon checks if a point lies within a polygon using ray-casting algorithm
func isPointInPolygon(longitude, latitude float64, polygon []Point) bool {
    point := Point{X: longitude, Y: latitude}
    inside := false
    j := len(polygon) - 1

    for i := 0; i < len(polygon); i++ {
        if ((polygon[i].Y > point.Y) != (polygon[j].Y > point.Y)) &&
            (point.X < (polygon[j].X-polygon[i].X)*(point.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X) {
            inside = !inside
        }
        j = i
    }

    return inside
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

    // Get all taxis and places
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
            // Log the polygon data being checked
            log.Printf("Checking against place %s with polygon: %v", 
                place.PlaceName, place.Polygon)
            
            // Unmarshal the polygon
            var polygon []Point
            polygonBytes, err := json.Marshal(place.Polygon)
            if err != nil {
                log.Printf("Error marshaling polygon for place %s: %v", place.PlaceName, err)
                continue
            }
            err = json.Unmarshal(polygonBytes, &polygon)
            if err != nil {
                log.Printf("Error unmarshaling polygon for place %s: %v", place.PlaceName, err)
                continue
            }

            // Check if taxi is within polygon
            if isPointInPolygon(taxi.Longitude, taxi.Latitude, polygon) {
                matched = true
                log.Printf("Taxi %s is within %s", taxi.TaxiID, place.PlaceName)
                
                // Insert or update mapping
                err := s.Repo.MappingRepo.InsertMapping(taxi.TaxiID, place.PlaceID)
                if err != nil {
                    log.Printf("Error inserting mapping: %v", err)
                    continue
                }

                // Update counter
                count, err := s.Repo.MappingRepo.GetCounter(taxi.TaxiID, place.PlaceID)
                if err != nil && err != sql.ErrNoRows {
                    log.Printf("Error getting counter: %v", err)
                    continue
                }

                if count == 0 {
                    err = s.Repo.MappingRepo.InsertCounter(taxi.TaxiID, place.PlaceID)
                } else {
                    err = s.Repo.MappingRepo.UpdateCounter(taxi.TaxiID, place.PlaceID)
                }
                if err != nil {
                    log.Printf("Error updating counter: %v", err)
                }
                break
            }
        }
        if !matched {
            log.Printf("Taxi %s does not match any place", taxi.TaxiID)
        }
    }
    log.Println("Mapping process completed")
}
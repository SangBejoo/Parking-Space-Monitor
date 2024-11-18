// internal/handlers/mapping.go
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/SangBejoo/parking-space-monitor/internal/scheduler"
)

// MappingHandler handles HTTP requests for Mapping operations.
type MappingHandler struct {
    Scheduler *scheduler.Scheduler
}

// TriggerMapping manually triggers the mapTaxiLocations function.
func (mh *MappingHandler) TriggerMapping(w http.ResponseWriter, r *http.Request) {
    go mh.Scheduler.MapTaxiLocations() // Run in a separate goroutine to prevent blocking
    fmt.Fprintf(w, "Mapping process triggered manually.")
}

// GetMapping retrieves current mappings with counters.
func (mh *MappingHandler) GetMapping(w http.ResponseWriter, r *http.Request) {
    query := `
        SELECT m.taxi_id, p.place_name, c.counter 
        FROM mapping m 
        JOIN places p ON m.place_id = p.place_id 
        JOIN counters c ON m.taxi_id = c.taxi_id AND m.place_id = c.place_id
    `
    rows, err := mh.Scheduler.Repo.DB.Query(query)
    if err != nil {
        http.Error(w, "Failed to query mappings: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var mappings []map[string]interface{}
    count := 0
    for rows.Next() {
        var taxiID, placeName string
        var counter int
        if err := rows.Scan(&taxiID, &placeName, &counter); err != nil {
            http.Error(w, "Failed to scan mapping: "+err.Error(), http.StatusInternalServerError)
            return
        }
        mappings = append(mappings, map[string]interface{}{
            "taxi_id": taxiID,
            "place":   placeName,
            "counter": counter,
        })
        count++
    }

    if err = rows.Err(); err != nil {
        http.Error(w, "Row iteration error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // If no mappings found, return an empty array instead of null
    w.Header().Set("Content-Type", "application/json")
    if count == 0 {
        json.NewEncoder(w).Encode([]map[string]interface{}{})
        return
    }
    json.NewEncoder(w).Encode(mappings)
}
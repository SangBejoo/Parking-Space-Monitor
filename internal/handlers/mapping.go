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
    mappings, err := mh.Scheduler.Repo.MappingRepo.GetAllMappings()
    if err != nil {
        http.Error(w, "Failed to retrieve mappings: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mappings)
}

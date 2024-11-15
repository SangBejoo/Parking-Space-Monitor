package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SangBejoo/parking-space-monitor/internal/scheduler"
	_ "github.com/lib/pq" // Required for postgres driver
)

func setupPostgresDB() (*sql.DB, error) {
	// Update connection string to use standard postgres format
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", // host
		5431,        // port
		"root",      // username
		"secret",    // password
		"subagiya1",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	// Test connection with timeout
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}

	// Create tables if they don't exist - using basic point checking instead of PostGIS
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS places (
            place_id SERIAL PRIMARY KEY,
            place_name VARCHAR(255) NOT NULL,
            min_lon FLOAT NOT NULL,
            min_lat FLOAT NOT NULL,
            max_lon FLOAT NOT NULL,
            max_lat FLOAT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS taxi (
            taxi_id VARCHAR(255) PRIMARY KEY,
            longitude FLOAT NOT NULL,
            latitude FLOAT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS mapping (
            taxi_id VARCHAR(255) REFERENCES taxi(taxi_id),
            place_id INTEGER REFERENCES places(place_id),
            PRIMARY KEY (taxi_id, place_id)
        );

        CREATE TABLE IF NOT EXISTS counters (
            taxi_id VARCHAR(255) REFERENCES taxi(taxi_id),
            place_id INTEGER REFERENCES places(place_id),
            counter INTEGER DEFAULT 0,
            last_counted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (taxi_id, place_id)
        );
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	// Clean up existing test data
	_, err = db.Exec(`
        TRUNCATE TABLE places, taxi, mapping, counters RESTART IDENTITY CASCADE
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to truncate tables: %v", err)
	}

	return db, nil
}

func insertTestData(db *sql.DB) error {
	// Insert 5 places with bounding box coordinates instead of polygons
	places := []struct {
		name   string
		minLon float64
		minLat float64
		maxLon float64
		maxLat float64
	}{
		{"Airport", 106.65, -6.13, 106.66, -6.12},
		{"Mall", 106.82, -6.20, 106.83, -6.19},
		{"Station", 106.84, -6.22, 106.85, -6.21},
		{"Hospital", 106.80, -6.18, 106.81, -6.17},
		{"Park", 106.83, -6.24, 106.84, -6.23},
	}

	for _, p := range places {
		_, err := db.Exec(`
            INSERT INTO places (place_name, min_lon, min_lat, max_lon, max_lat) 
            VALUES ($1, $2, $3, $4, $5)
        `, p.name, p.minLon, p.minLat, p.maxLon, p.maxLat)
		if err != nil {
			return fmt.Errorf("failed to insert place: %v", err)
		}
	}

	// Insert 50 taxis with realistic coordinates (Jakarta area)
	for i := 1; i <= 50; i++ {
		// Generate random coordinates within Jakarta
		lon := 106.7 + (float64(i) * 0.001) // Range: 106.7 - 106.75
		lat := -6.15 + (float64(i) * 0.001) // Range: -6.15 - -6.10

		_, err := db.Exec(`
            INSERT INTO taxi (taxi_id, longitude, latitude) 
            VALUES ($1, $2, $3)
        `, fmt.Sprintf("TEST-TAXI-%d", i), lon, lat)
		if err != nil {
			return fmt.Errorf("failed to insert taxi: %v", err)
		}
	}

	return nil
}

func TestRealPostgresMapping(t *testing.T) {
	db, err := setupPostgresDB()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	err = insertTestData(db)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	scheduler := &scheduler.Scheduler{Repo: &scheduler.Repository{DB: db}}
	handler := &MappingHandler{Scheduler: scheduler}

	// Trigger mapping
	reqTrigger, _ := http.NewRequest("POST", "/mapping/trigger", nil)
	rrTrigger := httptest.NewRecorder()
	handler.TriggerMapping(rrTrigger, reqTrigger)

	// Wait a bit for mapping to complete
	time.Sleep(2 * time.Second)

	// Get mapping results
	reqGet, _ := http.NewRequest("GET", "/mapping", nil)
	rrGet := httptest.NewRecorder()
	handler.GetMapping(rrGet, reqGet)

	if status := rrGet.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var mappings []map[string]interface{}
	if err := json.NewDecoder(rrGet.Body).Decode(&mappings); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify results
	t.Logf("Found %d mappings", len(mappings))
	for _, m := range mappings {
		t.Logf("Taxi %v is in place %v with counter %v",
			m["taxi_id"], m["place"], m["counter"])
	}
}

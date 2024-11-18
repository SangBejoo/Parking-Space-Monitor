
CREATE TABLE IF NOT EXISTS taxi_mapping (
    id SERIAL PRIMARY KEY,
    taxi_id VARCHAR(255) NOT NULL,
    place_id INTEGER NOT NULL,
    duration INTEGER DEFAULT 0,
    FOREIGN KEY (place_id) REFERENCES places(place_id)
);
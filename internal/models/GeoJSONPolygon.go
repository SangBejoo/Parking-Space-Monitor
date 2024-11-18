
package models

import "encoding/json"



// GeoJSONPolygon represents a GeoJSON polygon.

type GeoJSONPolygon struct {

    Type        string            `json:"type"`

    Coordinates [][][]json.Number `json:"coordinates"`

}



// Define other types here

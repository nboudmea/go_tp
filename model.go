package main

import "time"

// Modèle interne unifié — neutre, sans tag JSON/XML.
// On y mappe les deux sources (JSON et XML) après parsing.

type Station struct {
	ID           string
	Name         string
	CountryCode  string // ISO-2 ("FR", "ES", …)
	Coord        Coordinates
	Device       Device
	Observations []Observation
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
	Altitude  int // mètres
}

type Device struct {
	Model        string
	Manufacturer string
	InstalledOn  time.Time
}

type Observation struct {
	Timestamp     time.Time
	Temperature   float64 // °C
	Humidity      int     // %
	Pressure      float64 // hPa
	Wind          Wind
	Precipitation float64 // mm
	AirQuality    AirQuality
	Conditions    string
	Notes         *string
}

type Wind struct {
	Speed     float64 // km/h
	Direction int     // degrés 0-359
}

type AirQuality struct {
	PM25 float64
	PM10 float64
	NO2  float64
}

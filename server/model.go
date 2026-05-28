package main

import "time"

type Station struct {
	ID           string
	Name         string
	CountryCode  string
	Coord        Coordinates
	Device       Device
	Observations []Observation
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
	Altitude  int
}

type Device struct {
	Model        string
	Manufacturer string
	InstalledOn  time.Time
}

type Observation struct {
	Timestamp     time.Time
	Temperature   float64
	Humidity      int
	Pressure      float64
	Wind          Wind
	Precipitation float64
	AirQuality    AirQuality
	Conditions    string
	Notes         *string
}

type Wind struct {
	Speed     float64
	Direction int
}

type AirQuality struct {
	PM25 float64
	PM10 float64
	NO2  float64
}

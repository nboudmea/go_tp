package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Types intermédiaires privés — collent à la structure du fichier JSON.
// Ils ne fuient jamais hors de ce fichier.

type jsonDoc struct {
	Metadata struct {
		StationCount int `json:"station_count"`
	} `json:"metadata"`
	Stations []jsonStation `json:"stations"`
}

type jsonStation struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Country    string `json:"country"`     // "France", "Espagne", … (nom complet français)
	AltitudeM  int    `json:"altitude_m"`
	Location   struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Device struct {
		Type         string `json:"type"`
		Manufacturer string `json:"manufacturer"`
		InstalledOn  string `json:"installed_on"`
	} `json:"device"`
	Observations []jsonObs `json:"observations"`
}

type jsonObs struct {
	Timestamp          string  `json:"timestamp"`
	TemperatureCelsius float64 `json:"temperature_celsius"`
	HumidityPercent    int     `json:"humidity_percent"`
	PressureHpa        float64 `json:"pressure_hpa"`
	Wind               struct {
		SpeedKmh     float64 `json:"speed_kmh"`
		DirectionDeg int     `json:"direction_deg"`
	} `json:"wind"`
	PrecipitationMm float64 `json:"precipitation_mm"`
	AirQuality      struct {
		PM25 float64 `json:"pm25"`
		PM10 float64 `json:"pm10"`
		NO2  float64 `json:"no2"`
	} `json:"air_quality"`
	Conditions string  `json:"conditions"`
	Notes      *string `json:"notes"`
}

// Conversion nom complet français -> code ISO-2.
// Suffit pour le dataset fourni (14 pays).
var frenchCountryToISO = map[string]string{
	"France":    "FR",
	"Espagne":   "ES",
	"Portugal":  "PT",
	"Italie":    "IT",
	"Allemagne": "DE",
	"Belgique":  "BE",
	"Pays-Bas":  "NL",
	"Autriche":  "AT",
	"Suisse":    "CH",
	"Danemark":  "DK",
	"Suède":     "SE",
	"Norvège":   "NO",
	"Pologne":   "PL",
	"Tchéquie":  "CZ",
}

// LoadFromJSON lit le fichier puis le mappe vers le modèle interne.
func LoadFromJSON(path string) ([]Station, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lecture fichier JSON: %w", err)
	}
	var doc jsonDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("décodage JSON: %w", err)
	}
	stations := make([]Station, 0, len(doc.Stations))
	for _, js := range doc.Stations {
		st, err := jsonStationToStation(js)
		if err != nil {
			return nil, fmt.Errorf("station %s: %w", js.ID, err)
		}
		stations = append(stations, st)
	}
	return stations, nil
}

func jsonStationToStation(js jsonStation) (Station, error) {
	installed, err := time.Parse("2006-01-02", js.Device.InstalledOn)
	if err != nil {
		return Station{}, fmt.Errorf("installed_on invalide: %w", err)
	}
	iso, ok := frenchCountryToISO[js.Country]
	if !ok {
		return Station{}, fmt.Errorf("pays inconnu: %q", js.Country)
	}
	obs := make([]Observation, 0, len(js.Observations))
	for _, jo := range js.Observations {
		ts, err := time.Parse(time.RFC3339, jo.Timestamp)
		if err != nil {
			return Station{}, fmt.Errorf("timestamp invalide: %w", err)
		}
		obs = append(obs, Observation{
			Timestamp:   ts,
			Temperature: jo.TemperatureCelsius,
			Humidity:    jo.HumidityPercent,
			Pressure:    jo.PressureHpa,
			Wind: Wind{
				Speed:     jo.Wind.SpeedKmh,
				Direction: jo.Wind.DirectionDeg,
			},
			Precipitation: jo.PrecipitationMm,
			AirQuality: AirQuality{
				PM25: jo.AirQuality.PM25,
				PM10: jo.AirQuality.PM10,
				NO2:  jo.AirQuality.NO2,
			},
			Conditions: jo.Conditions,
			Notes:      jo.Notes,
		})
	}
	return Station{
		ID:          js.ID,
		Name:        js.Name,
		CountryCode: iso,
		Coord: Coordinates{
			Latitude:  js.Location.Latitude,
			Longitude: js.Location.Longitude,
			Altitude:  js.AltitudeM,
		},
		Device: Device{
			Model:        js.Device.Type,
			Manufacturer: js.Device.Manufacturer,
			InstalledOn:  installed,
		},
		Observations: obs,
	}, nil
}

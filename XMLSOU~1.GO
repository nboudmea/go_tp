package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Types intermédiaires privés — collent à la structure du fichier XML.

type xmlDoc struct {
	XMLName  xml.Name     `xml:"weather_dataset"`
	Version  string       `xml:"version,attr"`
	Source   string       `xml:"source,attr"`
	Stations []xmlStation `xml:"station"`
}

type xmlStation struct {
	ID         string `xml:"id,attr"`
	CountryISO string `xml:"country,attr"` // déjà en ISO-2
	Name       string `xml:"name"`
	Coordinates struct {
		Lat      float64 `xml:"lat,attr"`
		Lon      float64 `xml:"lon,attr"`
		Altitude int     `xml:"altitude,attr"`
	} `xml:"coordinates"`
	Hardware struct {
		Vendor string `xml:"vendor,attr"`
		Model  string `xml:"model,attr"`
		Since  string `xml:"since,attr"`
	} `xml:"hardware"`
	Observations struct {
		Count int      `xml:"count,attr"`
		Items []xmlObs `xml:"observation"`
	} `xml:"observations"`
}

type xmlObs struct {
	At       string       `xml:"at,attr"`
	Sky      string       `xml:"sky,attr"`
	Measures []xmlMeasure `xml:"measure"`
	Wind     struct {
		Speed     float64 `xml:"speed,attr"`
		Direction int     `xml:"direction,attr"`
	} `xml:"wind"`
	AirQuality struct {
		Pollutants []xmlPollutant `xml:"pollutant"`
	} `xml:"air_quality"`
	Note *string `xml:"note,omitempty"`
}

type xmlMeasure struct {
	Type  string `xml:"type,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:",chardata"`
}

type xmlPollutant struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func LoadFromXML(path string) ([]Station, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lecture fichier XML: %w", err)
	}
	var doc xmlDoc
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("décodage XML: %w", err)
	}
	stations := make([]Station, 0, len(doc.Stations))
	for _, xs := range doc.Stations {
		st, err := xmlStationToStation(xs)
		if err != nil {
			return nil, fmt.Errorf("station %s: %w", xs.ID, err)
		}
		stations = append(stations, st)
	}
	return stations, nil
}

func xmlStationToStation(xs xmlStation) (Station, error) {
	installed, err := time.Parse("2006-01-02", xs.Hardware.Since)
	if err != nil {
		return Station{}, fmt.Errorf("since invalide: %w", err)
	}
	obs := make([]Observation, 0, len(xs.Observations.Items))
	for _, xo := range xs.Observations.Items {
		ts, err := time.Parse(time.RFC3339, xo.At)
		if err != nil {
			return Station{}, fmt.Errorf("at invalide: %w", err)
		}
		o := Observation{
			Timestamp:  ts,
			Conditions: xo.Sky,
			Wind: Wind{
				Speed:     xo.Wind.Speed,
				Direction: xo.Wind.Direction,
			},
			Notes: xo.Note,
		}
		// Dispatch des mesures typées par attribut "type".
		for _, m := range xo.Measures {
			v, err := strconv.ParseFloat(m.Value, 64)
			if err != nil {
				return Station{}, fmt.Errorf("mesure %q invalide: %w", m.Type, err)
			}
			switch m.Type {
			case "temperature":
				o.Temperature = v
			case "humidity":
				o.Humidity = int(v)
			case "pressure":
				o.Pressure = v
			case "precipitation":
				o.Precipitation = v
			}
		}
		// Dispatch des polluants par attribut "name".
		for _, p := range xo.AirQuality.Pollutants {
			v, err := strconv.ParseFloat(p.Value, 64)
			if err != nil {
				return Station{}, fmt.Errorf("polluant %q invalide: %w", p.Name, err)
			}
			switch p.Name {
			case "PM2.5":
				o.AirQuality.PM25 = v
			case "PM10":
				o.AirQuality.PM10 = v
			case "NO2":
				o.AirQuality.NO2 = v
			}
		}
		obs = append(obs, o)
	}
	return Station{
		ID:          xs.ID,
		Name:        xs.Name,
		CountryCode: xs.CountryISO,
		Coord: Coordinates{
			Latitude:  xs.Coordinates.Lat,
			Longitude: xs.Coordinates.Lon,
			Altitude:  xs.Coordinates.Altitude,
		},
		Device: Device{
			Model:        xs.Hardware.Model,
			Manufacturer: xs.Hardware.Vendor,
			InstalledOn:  installed,
		},
		Observations: obs,
	}, nil
}

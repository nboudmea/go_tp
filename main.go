package main

import (
	"fmt"
	"os"
)

func main() {
	stationsJSON, err := LoadFromJSON("weather_data.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur JSON:", err)
		os.Exit(1)
	}
	stationsXML, err := LoadFromXML("weather_data.xml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur XML:", err)
		os.Exit(1)
	}

	// Comptes globaux
	totalObsJSON := 0
	for _, s := range stationsJSON {
		totalObsJSON += len(s.Observations)
	}
	totalObsXML := 0
	for _, s := range stationsXML {
		totalObsXML += len(s.Observations)
	}
	fmt.Printf("JSON : %d stations, %d observations\n", len(stationsJSON), totalObsJSON)
	fmt.Printf("XML  : %d stations, %d observations\n", len(stationsXML), totalObsXML)
	if len(stationsJSON) == len(stationsXML) && totalObsJSON == totalObsXML {
		fmt.Println("Cohérence : OK")
	} else {
		fmt.Println("Cohérence : KO")
	}

	// Station la plus ventée
	best, gust := MaxWindGust(stationsJSON)
	fmt.Printf("\nStation la plus ventée : %s (%.1f km/h)\n", best.ID, gust)

	// Moyenne température Bordeaux Mérignac
	if s, ok := FindByID(stationsJSON, "FR-BOR-001"); ok {
		fmt.Printf("Temp. moyenne %s : %.1f °C\n", s.Name, AvgTemperature(s))
	}

	// Compte par pays
	fmt.Printf("\nStations par pays : %v\n", CountByCountry(stationsJSON))
}

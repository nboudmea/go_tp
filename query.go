package main

// FilterByCountry renvoie les stations d'un pays donné (code ISO).
func FilterByCountry(stations []Station, iso string) []Station {
	out := make([]Station, 0)
	for _, s := range stations {
		if s.CountryCode == iso {
			out = append(out, s)
		}
	}
	return out
}

// AvgTemperature renvoie la température moyenne (°C) d'une station.
// Renvoie 0 si la station n'a aucune observation.
func AvgTemperature(s Station) float64 {
	if len(s.Observations) == 0 {
		return 0
	}
	var sum float64
	for _, o := range s.Observations {
		sum += o.Temperature
	}
	return sum / float64(len(s.Observations))
}

// MaxWindGust renvoie la station avec la rafale de vent la plus forte
// sur l'ensemble du dataset, avec la valeur correspondante.
func MaxWindGust(stations []Station) (Station, float64) {
	var best Station
	var max float64
	for _, s := range stations {
		for _, o := range s.Observations {
			if o.Wind.Speed > max {
				max = o.Wind.Speed
				best = s
			}
		}
	}
	return best, max
}

// CountByCountry renvoie le nombre de stations par code pays.
func CountByCountry(stations []Station) map[string]int {
	out := make(map[string]int)
	for _, s := range stations {
		out[s.CountryCode]++
	}
	return out
}

// FindByID renvoie la station correspondant à un ID, ou false.
// Utile pour le TP API CRUD.
func FindByID(stations []Station, id string) (Station, bool) {
	for _, s := range stations {
		if s.ID == id {
			return s, true
		}
	}
	return Station{}, false
}

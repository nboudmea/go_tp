package main

import (
	"encoding/json"
	"net/http"
)

type App struct{ store *Store }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (a *App) listStations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.All())
}

func (a *App) getStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	station, ok := a.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "station introuvable")
		return
	}
	writeJSON(w, http.StatusOK, station)
}

func (a *App) updateStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var station Station
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	station.ID = id
	if a.store.Has(id) {
		a.store.Put(station)
		writeJSON(w, http.StatusOK, station)
		return
	}
	a.store.Put(station)
	writeJSON(w, http.StatusCreated, station)
}

func (a *App) createStation(w http.ResponseWriter, r *http.Request) {
	var station Station
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	if a.store.Has(station.ID) {
		writeError(w, http.StatusConflict, "station déjà existante")
		return
	}
	a.store.Put(station)
	writeJSON(w, http.StatusCreated, station)
}

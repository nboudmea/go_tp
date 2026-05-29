package main

import (
	"encoding/json"
	"net/http"
)

type App struct{ store *Store }

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg, Code: code})
}

func (a *App) listStations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.All())
}

func (a *App) getStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	station, ok := a.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "station introuvable")
		return
	}
	writeJSON(w, http.StatusOK, station)
}

func (a *App) createStation(w http.ResponseWriter, r *http.Request) {
	var station Station
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "JSON invalide")
		return
	}
	if a.store.Has(station.ID) {
		writeError(w, http.StatusConflict, "ID_TAKEN", "station déjà existante")
		return
	}
	a.store.Put(station)
	writeJSON(w, http.StatusCreated, station)
}

func (a *App) updateStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var station Station
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "JSON invalide")
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

func (a *App) deleteStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !a.store.Delete(id) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "station introuvable")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) listObservations(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	station, ok := a.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "station introuvable")
		return
	}
	writeJSON(w, http.StatusOK, station.Observations)
}

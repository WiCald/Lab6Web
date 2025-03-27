package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Match representa un partido de La Liga, creo. No sé mucho de esto
type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"homeTeam"`
	AwayTeam    string `json:"awayTeam"`
	MatchDate   string `json:"matchDate"`
	Goals       int    `json:"goals"`
	YellowCards int    `json:"yellowCards"`
	RedCards    int    `json:"redcards"`
	ExtraTime   bool   `json:"extraTime"`
}

var (
	matches    = make(map[int]*Match)
	nextID     = 1
	matchesMux sync.Mutex
)

// Handlers de la API
func getMatches(w http.ResponseWriter, r *http.Request) {
	matchesMux.Lock()
	defer matchesMux.Unlock()

	list := make([]*Match, 0, len(matches))
	for _, match := range matches {
		list = append(list, match)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func getMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	matchesMux.Lock()
	defer matchesMux.Unlock()
	match, exists := matches[id]
	if !exists {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func createMatch(w http.ResponseWriter, r *http.Request) {
	var m Match
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	matchesMux.Lock()
	m.ID = nextID
	nextID++
	m.Goals = 0
	m.YellowCards = 0
	m.RedCards = 0
	m.ExtraTime = false
	matches[m.ID] = &m
	matchesMux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func updateMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var updated Match
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	matchesMux.Lock()
	defer matchesMux.Unlock()
	match, exists := matches[id]
	if !exists {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	match.HomeTeam = updated.HomeTeam
	match.AwayTeam = updated.AwayTeam
	match.MatchDate = updated.MatchDate
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func deleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	matchesMux.Lock()
	defer matchesMux.Unlock()
	if _, exists := matches[id]; !exists {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	delete(matches, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	// Agregar un handler opcional en la raíz para confirmar que el servidor está activo, por el amor de Dios, que funcione
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la API de La Liga Tracker"))
	})

	// Endpoints obligatorios con prefijo /api
	r.HandleFunc("/api/matches", getMatches).Methods("GET")
	r.HandleFunc("/api/matches/{id}", getMatch).Methods("GET")
	r.HandleFunc("/api/matches", createMatch).Methods("POST")
	r.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	r.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

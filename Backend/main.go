package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Middleware para habilitar CORS
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Encabezados CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Requested-With, Accept")
		// Si es una solicitud preflight (OPTIONS), responde inmediatamente
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Match representa un partido de La Liga.
type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"homeTeam"`
	AwayTeam    string `json:"awayTeam"`
	MatchDate   string `json:"matchDate"`
	Goals       int    `json:"goals"`
	YellowCards int    `json:"yellowCards"`
	RedCards    int    `json:"redCards"`
	ExtraTime   bool   `json:"extraTime"`
}

var (
	matches = make(map[int]*Match)
	nextID  = 1
	mu      sync.Mutex
)

// getMatches: GET /api/matches
func getMatches(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var list []*Match
	for _, match := range matches {
		list = append(list, match)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// getMatch: GET /api/matches/{id}
func getMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	mu.Lock()
	match, exists := matches[id]
	mu.Unlock()
	if !exists {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// createMatch: POST /api/matches
func createMatch(w http.ResponseWriter, r *http.Request) {
	var m Match
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	mu.Lock()
	m.ID = nextID
	nextID++
	// Inicializamos los contadores y bandera de tiempo extra.
	m.Goals = 0
	m.YellowCards = 0
	m.RedCards = 0
	m.ExtraTime = false
	matches[m.ID] = &m
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

// updateMatch: PUT /api/matches/{id}
func updateMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var updated Match
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	mu.Lock()
	match, exists := matches[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	match.HomeTeam = updated.HomeTeam
	match.AwayTeam = updated.AwayTeam
	match.MatchDate = updated.MatchDate
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// deleteMatch: DELETE /api/matches/{id}
func deleteMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	mu.Lock()
	if _, exists := matches[id]; !exists {
		mu.Unlock()
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	delete(matches, id)
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

// registerGoal: PATCH /api/matches/{id}/goals
func registerGoal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	mu.Lock()
	match, exists := matches[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	match.Goals++
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// registerYellowCard: PATCH /api/matches/{id}/yellowcards
func registerYellowCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	mu.Lock()
	match, exists := matches[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	match.YellowCards++
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// registerRedCard: PATCH /api/matches/{id}/redcards
func registerRedCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	mu.Lock()
	match, exists := matches[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	match.RedCards++
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// setExtraTime: PATCH /api/matches/{id}/extratime
func setExtraTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	mu.Lock()
	match, exists := matches[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	match.ExtraTime = true
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func main() {
	r := mux.NewRouter()

	// Aplica el middleware CORS globalmente
	r.Use(enableCors)

	// Handler global para solicitudes OPTIONS
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Requested-With, Accept")
		w.WriteHeader(http.StatusOK)
	})

	// Rutas
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la API de La Liga Tracker"))
	})
	r.HandleFunc("/api/matches", getMatches).Methods("GET")
	r.HandleFunc("/api/matches/{id}", getMatch).Methods("GET")
	r.HandleFunc("/api/matches", createMatch).Methods("POST")
	r.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	r.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")
	r.HandleFunc("/api/matches/{id}/goals", registerGoal).Methods("PATCH")
	r.HandleFunc("/api/matches/{id}/yellowcards", registerYellowCard).Methods("PATCH")
	r.HandleFunc("/api/matches/{id}/redcards", registerRedCard).Methods("PATCH")
	r.HandleFunc("/api/matches/{id}/extratime", setExtraTime).Methods("PATCH")

	log.Println("Servidor corriendo en el puerto 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

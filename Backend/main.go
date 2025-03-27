package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Middleware para habilitar CORS y registrar peticiones
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recibido %s request para %s", r.Method, r.URL)
		// Agrega encabezados CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Requested-With, Accept")
		// Si es preflight (OPTIONS), responder de inmediato
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

// getMatches devuelve todos los partidos.
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

// getMatch devuelve un partido por ID.
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

// createMatch crea un nuevo partido.
func createMatch(w http.ResponseWriter, r *http.Request) {
	var m Match
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	mu.Lock()
	m.ID = nextID
	nextID++
	// Valores iniciales
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

// updateMatch actualiza un partido existente.
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

// deleteMatch elimina un partido por su ID.
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

func main() {
	r := mux.NewRouter()
	r.Use(enableCors)

	// Agregar una ruta catch-all para OPTIONS (por si alguna ruta no es capturada)
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Handler para la raíz (opcional)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido a la API de La Liga Tracker"))
	})

	// Endpoints de la API
	r.HandleFunc("/api/matches", getMatches).Methods("GET")
	r.HandleFunc("/api/matches/{id}", getMatch).Methods("GET")
	r.HandleFunc("/api/matches", createMatch).Methods("POST")
	r.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	r.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")

	log.Println("Servidor corriendo en el puerto 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

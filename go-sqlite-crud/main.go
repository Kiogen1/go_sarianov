package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

type Game struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Studio string `json:"studio"`
	Year   int    `json:"year"`
	Sold   int    `json:"sold"`
}

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var db *sql.DB

// Middleware для CORS с логированием preflight-запросов
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы только с вашего фронтенда
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With")

		if r.Method == "OPTIONS" {
			log.Println("Preflight request:", r.Method, r.URL.Path, r.Header)
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Create
func createGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var cb Game
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("INSERT INTO str_games(name, studio, year, sold) VALUES(?, ?, ?, ?)", cb.Name, cb.Studio, cb.Year, cb.Sold)
	if err != nil {
		http.Error(w, "insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	cb.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cb)
}

// Read all
func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, name, studio, year, sold FROM str_games")
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	strGames := []Game{}
	for rows.Next() {
		var cb Game
		if err := rows.Scan(&cb.ID, &cb.Name, &cb.Studio, &cb.Year, &cb.Sold); err != nil {
			http.Error(w, "scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		strGames = append(strGames, cb)
	}
	json.NewEncoder(w).Encode(strGames)
}

// Read one
func getGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var cb Game
	err = db.QueryRow("SELECT id, name, studio, year, sold FROM str_games WHERE id = ?", id).
		Scan(&cb.ID, &cb.Name, &cb.Studio, &cb.Year, &cb.Sold)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cb)
}

// Update
func updateGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var cb Game
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE str_games SET name = ?, studio = ?, year = ?, sold = ? WHERE id = ?", cb.Name, cb.Studio, cb.Year, cb.Sold, id)
	if err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cb.ID = id
	json.NewEncoder(w).Encode(cb)
}

// Delete
func deleteGame(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM str_games WHERE id = ?", id)
	if err != nil {
		http.Error(w, "delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	var err error

	db, err = sql.Open("sqlite", "strGames.db")
	if err != nil {
		log.Fatalf("DB open: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	// Создание таблицы
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS str_games (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		studio TEXT NOT NULL,
		year INTEGER NOT NULL,
		sold INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
		);`)
	if err != nil {
		log.Fatalf("DB create: %v", err)
	}

	// Роутер
	router := mux.NewRouter()

	//Защищенные маршруты
	router.HandleFunc("/strGames", createGame).Methods("POST")
	router.HandleFunc("/strGames", getGames).Methods("GET")
	router.HandleFunc("/strGames/{id}", getGame).Methods("GET")
	router.HandleFunc("/strGames/{id}", updateGame).Methods("PUT")
	router.HandleFunc("/strGames/{id}", deleteGame).Methods("DELETE")

	// Запуск сервера
	log.Println("Сервер запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(router)))
}

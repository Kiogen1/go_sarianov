package main

import (
	"database/sql" // Кодирование и декодирование (Сериализация)
	"log"

	// HTTP-сервер, обработка запроса
	// Преобразование строки в число
	// Роутер
	"github.com/gorilla/mux"
	_ "modernc.org/sqlite" // Драйвер SQLite
)

type StrategicGames struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Studio string `json:"studio"`
	Year   int    `json:"year_founded"`
	Sold   int    `json:"copy_sold"`
}

var db *sql.DB // Переменая для соединения с базой данных

func main() {
	var err error

	// Открыть (создать) файл БД в текущей директории
	db, err = sql.Open("sqlite", "C:/GO/strategic_games.db")
	if err != nil {
		log.Fatalf("DB open: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1) // Ограничение на запись

	// Создать БД, если нету. Exec() выполнит запрос без возврата строк
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS strategic_games (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		studio TEXT NOT NULL,
		year INTEGER NOT NULL,
		sold INTEGER NOT NULL
	);`)
	if err != nil {
		log.Fatalf("DB open: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/carBrands", createCarBrand).Methods("POST")
	router.HandleFunc("/carBrands", getCarBrands).Methods("GET")
	router.HandleFunc("/carBrands/{id}", getCarBrand).Methods("GET")
	router.HandleFunc("/carBrands/{id}", updateCarBrand).Methods("PUT")
	router.HandleFunc("/carBrands/{id}", deleteCarBrand).Methods("DELETE")

	// Запуск сервера
	log.Println("Сервер запущен, порт :8080")

}

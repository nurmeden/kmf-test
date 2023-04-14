package controllers

import (
	"database/sql"
	"fmt"
	"kmf-test/internal/configs"
	"kmf-test/internal/handlers"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

var db *sql.DB

// https://nationalbank.kz/rss/get_rates.cfm?fdate=13.04.2023
// http://localhost:8080/currency/save/20.03.2023
func InitRoutes() {
	cfg, err := configs.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Open database connection
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Port, cfg.Database.Name)
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/currency/save/{date}", handlers.SaveCurrencyHandler).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", handlers.GetCurrencyHandler).Methods("GET")
	log.Println("Запуск веб-сервера на http://localhost:8080/ ")

	http.ListenAndServe(":8080", r)

}

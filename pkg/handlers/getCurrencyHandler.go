package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kmf-test/pkg/configs"
	"kmf-test/pkg/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func GetCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := configs.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	dateStr := vars["date"]
	code := vars["code"]

	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid date format.Correct format:02-01-2006", http.StatusBadRequest)
		return
	}
	var rows *sql.Rows

	if code != "" {
		rows, err = db.Query("SELECT * FROM R_CURRENCY WHERE A_DATE = @date AND CODE = @code", sql.Named("date", date), sql.Named("code", code))

	} else {
		rows, err = db.Query("SELECT * FROM R_CURRENCY WHERE A_DATE = @date", sql.Named("date", date))
	}
	if err != nil {
		http.Error(w, "Unable to get currency data", http.StatusInternalServerError)
		return
	}
	fmt.Printf("rows: %v\n", rows)
	defer rows.Close()
	var currencies []models.Currency
	for rows.Next() {
		var currency models.Currency
		err := rows.Scan(&currency.ID, &currency.Title, &currency.Code, &currency.Value, &currency.A_DATE)
		if err != nil {
			http.Error(w, "Unable to get currency data", http.StatusInternalServerError)
			return
		}
		currencies = append(currencies, currency)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Unable to get currency data", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(currencies)
	if err != nil {
		http.Error(w, "Unable to encode currency data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

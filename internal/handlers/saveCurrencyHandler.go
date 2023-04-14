package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"kmf-test/internal/configs"
	"kmf-test/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CurrencyResponse struct {
	Success bool `json:"success"`
}

var db *sql.DB

func SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем значение параметра date из URL
	vars := mux.Vars(r)
	date := vars["date"]

	// Формируем URL для запроса к API нац. банка
	url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", date)

	// Выполняем запрос к API нац. банка и получаем ответ в виде XML
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting data from API: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Парсим XML из ответа
	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var data struct {
		XMLName xml.Name `xml:"rates"`
		Date    string   `xml:"date"`
		Items   []struct {
			Fullname    string  `xml:"fullname"`
			Title       string  `xml:"title"`
			Description float64 `xml:"description"`
		} `xml:"item"`
	}

	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		log.Println("Error while unmarshaling XML data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Сохраняем данные в базу данных асинхронно
	go func() {
		// saveCurrencyToDB(db, date, data.Items)
		for _, item := range data.Items {
			err := saveCurrencyToDB(item, date)
			if err != nil {
				log.Println("Error while saving currency to DB:", err)
			}
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func saveCurrencyToDB(item models.Item, date string) error {
	// Подключаемся к базе данных
	cfg, err := configs.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;charset=utf8",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Сохраняем данные в таблицу R_CURRENCY
	query := "INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (@Title, @Code, @Value, @Date)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error while preparing query:", err.Error())
		return err
	}
	defer stmt.Close()

	fmt.Printf("item.Title: %v\n", item.Fullname)
	_, err = stmt.Exec(sql.Named("Title", item.Fullname), sql.Named("Code", item.Title), sql.Named("Value", item.Description), sql.Named("Date", date))
	if err != nil {
		log.Println("Error while saving currency to DB:", err.Error())
		return err
	}

	return nil
}

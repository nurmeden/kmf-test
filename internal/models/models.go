package models

import "time"

type Config struct {
	DatabaseConnectionString string `json:"databaseConnectionString"`
	ServerPort               int    `json:"serverPort"`
}

type Item struct {
	Fullname    string  `xml:"fullname"`
	Title       string  `xml:"title"`
	Description float64 `xml:"description"`
}

type Currency struct {
	ID     int        `json:"id"`
	Title  string     `json:"title"`
	Code   string     `json:"code"`
	Value  float64    `json:"value"`
	A_DATE *time.Time `json:"a_date"`
}

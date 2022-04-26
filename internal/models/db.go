package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB

func InitializeDB() {
	connStr := "user=" + viper.GetString("user") + " password=" + viper.GetString("password") + " dbname=" + viper.GetString("dbName") + "sslmode=" + viper.GetString("sslMode")

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
}

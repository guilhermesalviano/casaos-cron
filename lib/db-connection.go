package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	entities "google-flights-crawler/entities"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func DbConnection(cfg DBConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao configurar o banco: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar no MariaDB: %v", err)
	}

	log.Println("✅ Conectado ao MariaDB com sucesso!")

	return db
}

func InsertIntoDB(db *sql.DB, r *entities.SearchResult) error {
	_, err := db.Exec(
		"INSERT INTO flight_crawled (origin, destination, airline, stops, price, flightDate, searchDate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		r.Origin,
		r.Destination,
		r.BestFlights[0].Airline, 
		r.BestFlights[0].Stops, 
		r.BestFlights[0].Price, 
		r.Date,
		time.Now(),
	)
	return err
}

func CreateDatabaseConnection() *sql.DB {
	cfg := DBConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}
	return DbConnection(cfg)
}
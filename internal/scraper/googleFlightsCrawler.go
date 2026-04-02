package scraper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"google-flights-crawler/domain"
	"google-flights-crawler/internal/flights"
	"google-flights-crawler/internal/notify"
	db "google-flights-crawler/internal/store"
	"log"
	"os"
	"time"
)


func GoogleFlightsCrawler(params domain.SearchParams, output *string) {
	notify.Notify(fmt.Sprintf("Starting crawling google flights %s → %s on %s...\n", params.DepartureID, params.ArrivalID, params.OutboundDate))

	result, err := flights.ScrapeFlights(params)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	db := db.CreateDatabaseConnectionFactory()
	defer db.Close()

	log.Printf("✅ Search completed: %s → %s on %s. Best price: %.0f %s", result.Origin, result.Destination, result.Date, result.BestPrice, result.Currency)

	er := saveFlightsInDB(db, result)
	if er != nil {
		log.Printf("Error: Could not insert into database: %v", er)
	} else {
		log.Printf("Success: Search saved to database")
	}

	if output != nil && *output != "" {
		data, _ := json.MarshalIndent(result, "", "  ")
		if err := os.WriteFile(*output, data, 0644); err != nil {
			log.Printf("Warning: Could not write output file: %v\n", err)
		} else {
			log.Printf("Success: Results saved to %s", *output)
		}
	}
}

func saveFlightsInDB(db *sql.DB, r *domain.SearchResult) error {
	_, err := db.Exec(
		"INSERT INTO flight_crawled (origin, destination, airline, stops, price, flightDate, searchDate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		r.Origin,
		r.Destination,
		r.BestFlights[0].Airline, 
		r.BestFlights[0].Stops, 
		r.BestFlights[0].Price, 
		r.Date,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
}
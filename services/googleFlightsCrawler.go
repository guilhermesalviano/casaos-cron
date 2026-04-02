package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"google-flights-crawler/lib"
	"google-flights-crawler/notifier"
)


func GoogleFlightsCrawler(params lib.SearchParams, output *string) {
	notifier.Notify(fmt.Sprintf("Starting crawling google flights %s → %s on %s...\n", params.DepartureID, params.ArrivalID, params.OutboundDate))

	result, err := lib.ScrapeFlights(params)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	db := lib.CreateDatabaseConnection()
	defer db.Close()

	log.Printf("✅ Search completed: %s → %s on %s. Best price: %.0f %s", result.Origin, result.Destination, result.Date, result.BestPrice, result.Currency)

	er := lib.SaveFlightsInDB(db, result)
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
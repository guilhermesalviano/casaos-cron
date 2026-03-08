package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"

	entities "google-flights-crawler/entities"
	lib "google-flights-crawler/lib"
	notifier "google-flights-crawler/notifier"
	utils "google-flights-crawler/utils"

	_ "github.com/go-sql-driver/mysql"
)

type Scheduler struct {
	*gocron.Scheduler
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	path := os.Getenv("SCHEDULERS_FILE_PATH")

	local, _ := time.LoadLocation("America/Sao_Paulo")
	scheduler := &Scheduler{gocron.NewScheduler(local)}

	scheduler.scheduleFlightsCrawler(utils.LoadSchedulersFromCSV(path, "flights"), utils.GetApiKey())
	scheduler.scheduleAmazonWishlistCrawler(utils.LoadSchedulersFromCSV(path, "wishlists"))

	scheduler.StartBlocking()
}

func (scheduler *Scheduler) scheduleAmazonWishlistCrawler(wishlists []utils.SchedulersCsv) {
	for _, wishItem := range wishlists {
		log.Printf("📅 Schedule Amazon Wishlist Crawler (every %s at %s)",
			wishItem.Day, wishItem.Time)

		_, err := utils.ScheduleOnDay(scheduler.Scheduler, wishItem.Day).At(wishItem.Time).Do(func() {
			startWishlistAmazonCrawler()
		})

		if err != nil {
			notifier.Notify(fmt.Sprintf("Error scheduling job: %s", err))
			os.Exit(1)
		}
	}
}

func (scheduler *Scheduler) scheduleFlightsCrawler(flights []utils.SchedulersCsv, apiKey *string) {
	for _, flight := range flights {
		log.Printf("📅 Schedule: %s → %s on %s (every %s at %s)",
			flight.DepartureID, flight.ArrivalID, flight.OutboundDate, flight.Day, flight.Time)

		params := lib.SearchParams{
			APIKey:       *apiKey,
			DepartureID:  flight.DepartureID,
			ArrivalID:    flight.ArrivalID,
			OutboundDate: flight.OutboundDate,
			ReturnDate:   flight.ReturnDate,
			Adults:       flight.Adults,
			TravelClass:  flight.TravelClass,
			Stops:        flight.Stops,
			Currency:     flight.Currency,
			Language:     flight.Language,
			Country:      flight.Country,
		}

		_, err := utils.ScheduleOnDay(scheduler.Scheduler, flight.Day).At(flight.Time).Do(func() {
			startGoogleFlightsCrawler(params, nil)
		})

		if err != nil {
			notifier.Notify(fmt.Sprintf("Error scheduling job: %s", err))
			os.Exit(1)
		}
	}
}

func printResults(r *entities.SearchResult) {
	fmt.Printf("\n╔══════════════════════════════════════════════════════╗\n")
	fmt.Printf("║         Google Flights — Best Price Crawler          ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════╝\n")
	fmt.Printf("  Route     : %s → %s\n", r.Origin, r.Destination)
	fmt.Printf("  Outbound  : %s\n", r.Date)
	if r.ReturnDate != "" {
		fmt.Printf("  Return    : %s\n", r.ReturnDate)
	}
	fmt.Printf("  Searched  : %s\n", r.SearchedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("  ★ Best price: %.0f %s\n\n", r.BestPrice, r.Currency)

	printSection := func(label string, flights []entities.Flight) {
		if len(flights) == 0 {
			return
		}
		// Sort by price
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price < flights[j].Price
		})
		fmt.Printf("── %s (%d results) ──────────────────────────────\n", label, len(flights))
		fmt.Printf("  %-22s %-8s %-8s %-10s %-6s %s\n",
			"Airline", "Dep", "Arr", "Duration", "Stops", "Price")
		fmt.Printf("  %s\n", "─────────────────────────────────────────────────────────")
		for _, f := range flights {
			dur := fmt.Sprintf("%dh%02dm", f.Duration/60, f.Duration%60)
			stops := "nonstop"
			if f.Stops == 1 {
				stops = "1 stop"
			} else if f.Stops > 1 {
				stops = fmt.Sprintf("%d stops", f.Stops)
			}
			fmt.Printf("  %-22s %-8s %-8s %-10s %-6s %.0f %s\n",
				utils.Truncate(f.Airline, 22),
				utils.TimeOnly(f.Departure),
				utils.TimeOnly(f.Arrival),
				dur, stops,
				f.Price, f.Currency,
			)
		}
		fmt.Println()
	}

	printSection("Best Flights", r.BestFlights)
	printSection("Other Flights", r.OtherFlights)
}

func startWishlistAmazonCrawler() {
	notifier.Notify("Starting crawling amazon wishlist...")
	results, err := lib.ScrapeAmazonWishlist()
	if err != nil {
		notifier.Notify("Error scraping Amazon wishlist: " + err.Error())
		os.Exit(1)
	}
	notifier.Notify(fmt.Sprintf("Amazon wishlist scraped successfully: %d items found", len(results)))

	db := lib.CreateDatabaseConnection()
	defer db.Close()

	var er error
	for index, result := range results {
		er = lib.SaveWishlistAmazonPricesInDB(db, &result)
		if er == nil {
			log.Printf("💾 Index %d saved at database.", index)
		}
	}
	if er != nil {
		notifier.Notify(fmt.Sprintf("⚠️ Could not insert into database: %v\n", er))
	}
}

func startGoogleFlightsCrawler(params lib.SearchParams, output *string) {
	notifier.Notify(fmt.Sprintf("🔍 Searching flights %s → %s on %s...\n", params.DepartureID, params.ArrivalID, params.OutboundDate))

	result, err := lib.ScrapeFlights(params)
	if err != nil {
		notifier.Notify(fmt.Sprintf("❌  Error: %v\n", err))
		os.Exit(1)
	}

	db := lib.CreateDatabaseConnection()
	defer db.Close()

	printResults(result)

	notifier.Notify(fmt.Sprintf("✅ Search completed: %s → %s on %s. Best price: %.0f %s", result.Origin, result.Destination, result.Date, result.BestPrice, result.Currency))

	er := lib.SaveFlightsInDB(db, result)
	if er != nil {
		notifier.Notify(fmt.Sprintf("⚠️ Could not insert into database: %v\n", er))
	} else {
		notifier.Notify("💾 Search saved to database")
	}

	if output != nil && *output != "" {
		data, _ := json.MarshalIndent(result, "", "  ")
		if err := os.WriteFile(*output, data, 0644); err != nil {
			notifier.Notify(fmt.Sprintf("⚠️  Could not write output file: %v\n", err))
		} else {
			notifier.Notify(fmt.Sprintf("💾 Results saved to %s", *output))
		}
	}
}

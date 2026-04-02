package main

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"

	"google-flights-crawler/domain"
	"google-flights-crawler/internal/ai"
	"google-flights-crawler/internal/notify"
	"google-flights-crawler/internal/scraper"
	"google-flights-crawler/utils"
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

func aiAnalysis(prompt string) string {
    analysis, err := ai.AnalyzeWithGemini(prompt)
    if err != nil {
        log.Printf("❌ Gemini Error: %v\n", err)
        notify.Notify("❌ Gemini Error: " + err.Error())
    }
	return analysis
}

func (scheduler *Scheduler) scheduleAmazonWishlistCrawler(wishlists []utils.SchedulersCsv) {
	for _, wishItem := range wishlists {
		log.Printf("📅 Schedule Amazon Wishlist Crawler (every %s at %s)",
			wishItem.Day, wishItem.Time)

		_, err := utils.ScheduleOnDay(scheduler.Scheduler, wishItem.Day).At(wishItem.Time).Do(func() {
			scraper.AmazonWishlistCrawler()
		})

		if err != nil {
			notify.Notify(fmt.Sprintf("Error scheduling job: %s", err))
			os.Exit(1)
		}
	}
}

func (scheduler *Scheduler) scheduleFlightsCrawler(flights []utils.SchedulersCsv, apiKey *string) {
	for _, flight := range flights {
		log.Printf("📅 Schedule: %s → %s on %s (every %s at %s)",
			flight.DepartureID, flight.ArrivalID, flight.OutboundDate, flight.Day, flight.Time)

		params := domain.SearchParams{
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
			scraper.GoogleFlightsCrawler(params, nil)
		})

		if err != nil {
			notify.Notify(fmt.Sprintf("Error scheduling job: %s", err))
			os.Exit(1)
		}
	}
}

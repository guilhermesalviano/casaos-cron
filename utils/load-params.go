package utils

import (
	"flag"
	"google-flights-crawler/notifier"
	"os"
	"time"
)

type SearchParams struct {
	APIKey       string
	DepartureID  string // IATA code, e.g. "GRU"
	ArrivalID    string // IATA code, e.g. "JFK"
	OutboundDate string // YYYY-MM-DD
	ReturnDate   string // YYYY-MM-DD (empty = one-way)
	Adults       int
	TravelClass  int // 1=Economy 2=Premium 3=Business 4=First
	Stops        int // 0=Any 1=Nonstop 2=1stop 3=2stops
	Currency     string // e.g. "BRL", "USD"
	Language     string // e.g. "pt", "en"
	Country      string // e.g. "br", "us"
}

func GetApiKey() *string {
	if key := os.Getenv("SERPAPI_KEY"); key != "" {
		return &key
	}

	key := flag.String("key", "", "SerpApi API key")
	flag.Parse()

	if *key == "" {
		notifier.Notify("SERPAPI_KEY não definida. Use a env var ou o flag -key")
	}

	return key
}

func GetFlagsValuesOld() (SearchParams, *string) {
	apiKey := flag.String("key", os.Getenv("SERPAPI_KEY"), "SerpApi API key (or set SERPAPI_KEY env var)")
	from := flag.String("from", "GRU", "Departure IATA code (e.g. GRU, JFK, LHR)")
	to := flag.String("to", "JFK", "Arrival IATA code (e.g. JFK, GRU, CDG)")
	outbound := flag.String("date", time.Now().AddDate(0, 1, 0).Format("2006-01-02"), "Outbound date YYYY-MM-DD")
	returnDate := flag.String("return", "", "Return date YYYY-MM-DD (empty = one-way)")
	adults := flag.Int("adults", 1, "Number of adult passengers")
	class := flag.Int("class", 1, "Travel class: 1=Economy 2=Premium Economy 3=Business 4=First")
	stops := flag.Int("stops", 0, "Max stops: 0=Any 1=Nonstop 2=1stop 3=2stops")
	currency := flag.String("currency", "BRL", "Currency code (e.g. BRL, USD, EUR)")
	lang := flag.String("lang", "pt", "Language code (e.g. pt, en)")
	country := flag.String("country", "br", "Country code (e.g. br, us)")
	output := flag.String("output", "", "Save results to JSON file (optional)")

	params := SearchParams{
		APIKey:       *apiKey,
		DepartureID:  *from,
		ArrivalID:    *to,
		OutboundDate: *outbound,
		ReturnDate:   *returnDate,
		Adults:       *adults,
		TravelClass:  *class,
		Stops:        *stops,
		Currency:     *currency,
		Language:     *lang,
		Country:      *country,
	}

	return params, output
}
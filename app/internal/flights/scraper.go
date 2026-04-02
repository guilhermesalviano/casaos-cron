package flights

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"cron-to-casaos/domain"
)

func BuildQuery(p domain.SearchParams) url.Values {
	q := url.Values{}
	q.Set("engine", "google_flights")
	q.Set("api_key", p.APIKey)
	q.Set("departure_id", p.DepartureID)
	q.Set("arrival_id", p.ArrivalID)
	q.Set("outbound_date", p.OutboundDate)
	q.Set("adults", strconv.Itoa(p.Adults))
	q.Set("travel_class", strconv.Itoa(p.TravelClass))
	q.Set("stops", strconv.Itoa(p.Stops))
	q.Set("currency", p.Currency)
	q.Set("hl", p.Language)
	q.Set("gl", p.Country)

	if p.ReturnDate != "" {
		q.Set("return_date", p.ReturnDate)
		q.Set("type", "1")
	} else {
		q.Set("type", "2")
	}
	return q
}

func ScrapeFlights(p domain.SearchParams) (*domain.SearchResult, error) {
	const SERPAPIBASE = "https://serpapi.com/search"
	reqURL := fmt.Sprintf("%s?%s", SERPAPIBASE, BuildQuery(p).Encode())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var raw domain.SerpResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	if raw.Error != "" {
		return nil, fmt.Errorf("API error: %s", raw.Error)
	}

	result := &domain.SearchResult{
		SearchedAt:  time.Now().UTC(),
		Origin:      p.DepartureID,
		Destination: p.ArrivalID,
		Date:        p.OutboundDate,
		ReturnDate:  p.ReturnDate,
		Currency:    p.Currency,
	}

	parse := func(raw []domain.SerpFlight) []domain.Flight {
		var out []domain.Flight
		for _, sf := range raw {
			airline := sf.Airline
			flightNum := ""
			depTime := ""
			arrTime := ""
			if len(sf.Flights) > 0 {
				if airline == "" {
					airline = sf.Flights[0].Airline
				}
				flightNum = sf.Flights[0].FlightNumber
				depTime = sf.Flights[0].DepartureAirport.Time
				arrTime = sf.Flights[len(sf.Flights)-1].ArrivalAirport.Time
			}
			out = append(out, domain.Flight{
				Airline:      airline,
				FlightNumber: flightNum,
				Departure:    depTime,
				Arrival:      arrTime,
				Duration:     sf.TotalDuration,
				Stops:        len(sf.Layovers),
				Price:        float64(sf.Price),
				Currency:     p.Currency,
				CarbonEmitKg: sf.CarbonEmissions.ThisFlightKg,
			})
		}
		return out
	}

	result.BestFlights = parse(raw.BestFlights)
	result.OtherFlights = parse(raw.OtherFlights)

	// Find overall best (lowest) price
	all := append(result.BestFlights, result.OtherFlights...)
	if len(all) > 0 {
		best := all[0].Price
		for _, f := range all[1:] {
			if f.Price > 0 && f.Price < best {
				best = f.Price
			}
		}
		result.BestPrice = best
	}

	return result, nil
}

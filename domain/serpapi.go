package domain

type serpLayover struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Duration int    `json:"duration"`
}

type SerpFlight struct {
	Flights []struct {
		DepartureAirport struct{ Time string `json:"time"` } `json:"departure_airport"`
		ArrivalAirport   struct{ Time string `json:"time"` } `json:"arrival_airport"`
		Duration         int    `json:"duration"`
		Airline          string `json:"airline"`
		FlightNumber     string `json:"flight_number"`
	} `json:"flights"`
	Layovers        []serpLayover `json:"layovers"`
	TotalDuration   int `json:"total_duration"`
	CarbonEmissions struct {
		ThisFlightKg int `json:"this_flight"`
	} `json:"carbon_emissions"`
	Price   int    `json:"price"`
	Airline string `json:"airline,omitempty"`
}

type SerpResponse struct {
	BestFlights  []SerpFlight `json:"best_flights"`
	OtherFlights []SerpFlight `json:"other_flights"`
	Error        string       `json:"error,omitempty"`
}
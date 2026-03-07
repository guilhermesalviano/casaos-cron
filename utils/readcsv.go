package utils

import (
	"encoding/csv"
	"os"
	"strconv"

	notifier "google-flights-crawler/notifier"
)

type SchedulersCsv struct {
	DepartureID   string
	ArrivalID     string
	OutboundDate  string
	ReturnDate    string
	Adults        int
	TravelClass   int
	Stops         int
	Currency      string
	Language      string
	Country       string
	SchedulerType string
	Day           string
	Time          string
}

func LoadSearchParams(filePath string, SchedulerType string) ([]SchedulersCsv, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var allParams []SchedulersCsv

	for _, row := range records {
		adults, _ := strconv.Atoi(row[4])
		class, _ := strconv.Atoi(row[5])
		stops, _ := strconv.Atoi(row[6])
		SchedulerTypeCsv := row[0]
		
		if (SchedulerTypeCsv != SchedulerType) {
			continue;
		}

		var scheduler SchedulersCsv

		switch SchedulerTypeCsv {
		case "flights":
			scheduler = SchedulersCsv{
				DepartureID:  row[1],
				ArrivalID:    row[2],
				OutboundDate: row[3],
				ReturnDate:   row[4],
				Adults:       adults,
				TravelClass:  class,
				Stops:        stops,
				Currency:     row[8],
				Language:     row[9],
				Country:      row[10],
				Day:          row[11],
				Time:         row[12],
			}
		case "wishlists":
			scheduler = SchedulersCsv{
				Day:           row[11],
				Time:          row[12],
			}
		default:
			continue
		}

		allParams = append(allParams, scheduler)
	}

	return allParams, nil
}

func LoadSchedulersFromCSV(filePath string, SchedulerType string) ([]SchedulersCsv) {
	schedulers, err := LoadSearchParams(filePath, SchedulerType)
	if err != nil {
		schedulers, err = LoadSearchParams("./schedulers.csv", SchedulerType)
		if err != nil {
			notifier.Notify("Error loading search params: " + err.Error())
			os.Exit(1)
		}
	}
	return schedulers
}
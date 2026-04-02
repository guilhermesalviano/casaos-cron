package scraper

import "context"

type RawResult struct {
}

type Scraper interface {
	Scrape(ctx context.Context) ([]RawResult, error)
}

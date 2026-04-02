package services

import (
	"log"
	"google-flights-crawler/lib"
	"google-flights-crawler/notifier"
)


func AmazonWishlistCrawler() {
	notifier.Notify("Starting crawling amazon wishlist...")
	results, err := lib.ScrapeAmazonWishlist()
	if err != nil {
		log.Printf("Error: Failed in scraping Amazon wishlist: %v", err)
	}
	log.Printf("Amazon wishlist scraped successfully: %d items found", len(results))

	db := lib.CreateDatabaseConnection()
	defer db.Close()

	var er error
	for index, result := range results {
		er = lib.SaveWishlistAmazonPricesInDB(db, &result)
		if er == nil {
			log.Printf("Success: Index %d saved at database.", index)
		}
	}
	if er != nil {
		log.Printf("Error: Could not insert into database: %v\n", er)
	}
}


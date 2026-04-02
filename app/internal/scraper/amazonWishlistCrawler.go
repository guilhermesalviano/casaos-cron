package scraper

import (
	"database/sql"
	"log"
	"time"
	"cron-to-casaos/domain"
	"cron-to-casaos/internal/notify"
	db "cron-to-casaos/internal/store"
	"cron-to-casaos/internal/wishlist"
)

func AmazonWishlistCrawler() {
	notify.Notify("Starting crawling amazon wishlist...")
	results, err := wishlist.ScrapeAmazonWishlist()
	if err != nil {
		log.Printf("Error: Failed in scraping Amazon wishlist: %v", err)
	}
	log.Printf("Amazon wishlist scraped successfully: %d items found", len(results))

	db := db.CreateDatabaseConnectionFactory()
	defer db.Close()

	var er error
	for index, result := range results {
		er = saveWishlistAmazonPricesInDB(db, &result)
		if er == nil {
			log.Printf("Success: Index %d saved at database.", index)
		}
	}
	if er != nil {
		log.Printf("Error: Could not insert into database: %v\n", er)
	}
}

func saveWishlistAmazonPricesInDB(db *sql.DB, r *domain.WishlistItem) error {
	_, err := db.Exec(
		"INSERT INTO wishlist_amazon (title, price, link, search_date) VALUES (?, ?, ?, ?)",
		r.Title,
		r.Price,
		r.Link,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
}

package lib

import (
	"fmt"
	"log"
	"os"
	"strings"
	
	"github.com/gocolly/colly/v2"
	"google-flights-crawler/entities"
)

func ScrapeAmazonWishlist() ([]entities.WishlistItem, error) {
	var items []entities.WishlistItem

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		colly.AllowedDomains("www.amazon.com.br", "amazon.com.br"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
	})

	c.OnHTML("li.g-item-sortable", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a[id^='itemName_']"))
		link := e.ChildAttr("a[id^='itemName_']", "href")
		price := strings.TrimSpace(e.ChildText(".a-price .a-offscreen"))
		if price == "" {
			price = strings.TrimSpace(e.ChildText("span[id^='itemPrice_']"))
		}

		if title != "" {
			items = append(items, entities.WishlistItem{
				Title: title,
				Price: price,
				Link:  e.Request.AbsoluteURL(link),
			})
		}
	})

	c.OnHTML("a.wl-see-more", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if nextPage != "" {
			nextPageURL := e.Request.AbsoluteURL(nextPage)
			e.Request.Visit(nextPageURL)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error on %s: %v (Status: %d)", r.Request.URL, err, r.StatusCode)
	})

	wishlistID := os.Getenv("WISHLIST_ID")
	if wishlistID == "" {
		return nil, fmt.Errorf("WISHLIST_ID environment variable is required")
	}

	baseURL := fmt.Sprintf("https://www.amazon.com.br/hz/wishlist/ls/%s?_encoding=UTF8&sort=price-asc&filter=unpurchased", wishlistID)
	
	err := c.Visit(baseURL)
	if err != nil {
		return nil, err
	}
	c.Wait()

	return items, nil
}
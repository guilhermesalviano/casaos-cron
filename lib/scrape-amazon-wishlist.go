package lib

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type WishlistItem struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func ScrapeAmazonWishlist() ([]WishlistItem, error) {
	items := []WishlistItem{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Sec-Ch-Ua", `"Chromium";v="122", "Not(A:Bit";v="24", "Google Chrome";v="122"`)
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	c.OnHTML("li.g-item-sortable", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a[id^='itemName_']"))
		price := strings.TrimSpace(e.ChildText(".a-price .a-offscreen"))
		if price == "" {
			price = strings.TrimSpace(e.ChildText("span[id^='itemPrice_']"))
		}

		if title != "" {
			items = append(items, WishlistItem{Title: title, Price: price})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\nError:%v", r.Request.URL, r.StatusCode, err)
	})

	wishlistID := os.Getenv("WISHLIST_ID")
	if wishlistID == "" {
		log.Fatal("WISHLIST_ID not defined.")
	}
	c.Visit(fmt.Sprintf("https://www.amazon.com.br/hz/wishlist/ls/%s?_encoding=UTF&sort=price-asc8&ref_=wl_share", wishlistID))

	return items, nil
}
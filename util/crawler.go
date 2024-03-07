package crawler

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func crawler() {
	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("div", func(e *colly.HTMLElement) {

		if e.Attr("class") == "cb-col-100 cb-col cb-schdl cb-billing-plans-text" {

			fmt.Printf("Link found: %s\n", e.Text)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.cricbuzz.com/live-cricket-scores/87922/afg-vs-ire-1st-odi-afghanistan-v-ireland-in-uae-2024")
}

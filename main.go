package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nexidian/gocliselect"
)

type Matches struct {
	Index int
	Title string
	Link  string
}

func main() {
	c := colly.NewCollector()
	matches := []Matches{}

	c.OnHTML("div", func(e *colly.HTMLElement) {

		if e.Attr("class") == "cb-col-100 cb-col cb-schdl cb-billing-plans-text" {
			match := Matches{
				Index: len(matches),
				Title: ParseTitle(e.Text),
				Link:  e.ChildAttr("a", "href"),
			}
			matches = append(matches, match)
		}
	})

	c.Visit("https://www.cricbuzz.com/cricket-match/live-scores")

	fmt.Printf("Matches: %v\n", matches[0])

	menu := gocliselect.NewMenu("Choose the match")

	for _, match := range matches {
		menu.AddItem(match.Title, match.Link)
	}

	choice := menu.Display()

	fmt.Printf("Choice: %s\n", choice)
}

func ParseTitle(title string) string {
	title = strings.Split(title, ",")[0]
	return title
}

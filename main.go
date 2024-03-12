package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nexidian/gocliselect"
)

type Matches struct {
	Index int
	Title string
	Link  string
}

type Batsman struct {
	Name  string
	Runs  string
	Balls string
}

type Bowler struct {
	Name    string
	Runs    string
	Wickets string
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

	menu := gocliselect.NewMenu("Choose the match")

	for _, match := range matches {
		menu.AddItem(match.Title, match.Link)
	}

	choice := menu.Display()
	fmt.Println()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				score, currentBats, currentBowls := GetLiveScore(choice)
				fmt.Print("\033[2J")
				fmt.Printf("\r%s\n", score)
				fmt.Println("\nBatsman")
				for idx, batsman := range currentBats {
					star := ""
					if idx == 0 {
						star = "*"
					}
					fmt.Printf("%s%s %s(%s)\n", batsman.Name, star, batsman.Runs, batsman.Balls)
				}
				fmt.Println("\nBowler")
				for _, bowler := range currentBowls {
					fmt.Printf("%s %s-%s\n", bowler.Name, bowler.Wickets, bowler.Runs)
				}
			}
		}
	}()
	<-make(chan struct{})
}

func ParseTitle(title string) string {
	title = strings.Split(title, ",")[0]
	return title
}

func GetLiveScore(link string) (string, []Batsman, []Bowler) {
	score := ""
	currentBats := []Batsman{}
	flag := 0
	currentBowls := []Bowler{}
	c := colly.NewCollector()
	c.OnHTML("div", func(e *colly.HTMLElement) {

		if e.Attr("class") == "cb-col-100 cb-col cb-col-scores" || e.Attr("class") == "cb-col cb-col-100 cb-col-scores" {
			score = e.Text
		}
	})

	c.OnHTML("div.cb-min-inf.cb-col-100", func(e *colly.HTMLElement) {
		if flag == 0 {
			e.ForEach("div.cb-col.cb-col-100.cb-min-itm-rw", func(i int, e *colly.HTMLElement) {
				batsman := Batsman{
					Name:  e.ChildText("a"),
					Runs:  e.ChildText("div:nth-of-type(2)"),
					Balls: e.ChildText("div:nth-of-type(3)"),
				}
				currentBats = append(currentBats, batsman)
			})
			flag++
		} else {
			e.ForEach("div.cb-col.cb-col-100.cb-min-itm-rw", func(i int, e *colly.HTMLElement) {
				bowler := Bowler{
					Name:    e.ChildText("a"),
					Runs:    e.ChildText("div:nth-of-type(4)"),
					Wickets: e.ChildText("div:nth-of-type(5)"),
				}
				currentBowls = append(currentBowls, bowler)
			})
		}
	})

	c.Visit("https://www.cricbuzz.com" + link)
	return score, currentBats, currentBowls
}

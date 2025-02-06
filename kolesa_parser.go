package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Car struct {
	Title       string
	Description string
	Price       string
	MonthlyPay  string
	Date        string
	City        string
	Views       string
}

func main() {
	numCars := flag.Int("count", 5, "Number of cars to parse")
	usageType := flag.String("used", "all", "Usage type: all, new, used")
	state := flag.String("state", "", "Car state: cleared, crashed")
	flag.Parse()

	var baseURL string

	switch *usageType {
	case "new":
		baseURL = "https://kolesa.kz/cars/novye-avtomobili/"
	case "used":
		baseURL = "https://kolesa.kz/cars/avtomobili-s-probegom/"
	default:
		baseURL = "https://kolesa.kz/cars/"
	}

	if *state != "" {
		switch *state {
		case "cleared":
			baseURL += "?auto-custom=2"
		case "crashed":
			if strings.Contains(baseURL, "?") {
				baseURL += "&need-repair=1"
			} else {
				baseURL += "?need-repair=1"
			}
		}
	}

	cars := []Car{}
	page := 1

	for len(cars) < *numCars {
		url := baseURL
		if page > 1 {
			if strings.Contains(baseURL, "?") {
				url += "&page=" + strconv.Itoa(page)
			} else {
				url += "?page=" + strconv.Itoa(page)
			}
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to fetch URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("HTTP request failed with status: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalf("Failed to parse HTML: %v", err)
		}

		doc.Find("div.a-list div.a-card__info").Each(func(i int, s *goquery.Selection) {
			if len(cars) >= *numCars {
				return
			}

			title := strings.TrimSpace(s.Find("h5.a-card__title a").Text())
			description := strings.TrimSpace(s.Find("p.a-card__description").Text())
			price := strings.TrimSpace(s.Find("span.a-card__price").Text())
			monthlyPay := strings.TrimSpace(s.Find("div.month-payment__amount").Text())
			date := strings.TrimSpace(s.Find("div.a-card__footer span.a-card__param--date").Text())
			city := strings.TrimSpace(s.Find("div.a-card__footer span[data-test=region]").Text())
			views := strings.TrimSpace(s.Find("span.a-card__views").Text())

			car := Car{
				Title:       title,
				Description: description,
				Price:       price,
				MonthlyPay:  monthlyPay,
				Date:        date,
				City:        city,
				Views:       views,
			}
			cars = append(cars, car)
		})

		if len(cars) < *numCars {
			page++
		} else {
			break
		}
	}

	for i, car := range cars {
		fmt.Printf("Car %d:\n", i+1)
		fmt.Printf("  Title: %s\n", car.Title)
		fmt.Printf("  Description: %s\n", car.Description)
		fmt.Printf("  Price: %s\n", car.Price)
		fmt.Printf("  Monthly Pay: %s\n", car.MonthlyPay)
		fmt.Printf("  Date: %s\n", car.Date)
		fmt.Printf("  City: %s\n", car.City)
		fmt.Printf("  Views: %s\n\n", car.Views)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Engine struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Price    string `json:"price,optional"`
	Shipping string `json:"shipping,optional"`
	Img      string `json:"img,optional"`
	Grade    string `json:"grade,optional"`
}

var Engines []Engine

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func minmax(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func main() {
	// 2012 Challenger - 2C3CDYCJ5CH127691
	// 2019 Cadillac - 1G6AR5SX2K0139697
	start := time.Now()
	vin := "2C3CDYCJ5CH127691"
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/%s", vin))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err1 := json.Unmarshal(body, &Engines)
	if err1 != nil {
		log.Fatal(err1)
	}

	var enginePrices []float64

	for _, engine := range Engines {
		if engine.Price != "" {
			engine.Price = trimFirstRune(engine.Price)
			engine.Price = strings.Replace(engine.Price, ",", "", -1)
		}
		if _, err := strconv.ParseFloat(engine.Price, 64); err == nil {
			f, err := strconv.ParseFloat(engine.Price, 64)
			if err != nil {
				log.Fatal(err)
			}
			enginePrices = append(enginePrices, f)
		}
	}

	total := 0.0

	for _, price := range enginePrices {
		total = total + price
	}

	avg := total / float64(len(enginePrices))

	min, max := minmax(enginePrices)

	fmt.Println("Total Number of Engines: ", len(Engines))
	fmt.Println("Engines with Prices: ", len(enginePrices))
	fmt.Printf("Lowest Priced Engine: $%v \n", min)
	fmt.Printf("Highest Priced Engine: $%v \n", max)
	fmt.Printf("Average Engine Price for vehicle with VIN (%s): $%f", vin, avg)

	duration := time.Since(start)
	fmt.Println("\nTime taken: ", duration)

}

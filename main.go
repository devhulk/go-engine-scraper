package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

// Part - Engine, Engine Cover, Coolant Pump
type Part struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Price    string `json:"price,optional"`
	Shipping string `json:"shipping,optional"`
	Img      string `json:"img,optional"`
	Grade    string `json:"grade,optional"`
}

var parts []Part

func getParts(url_suffix string) []Part {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)

	c.OnHTML("div[class=individualPartHolder]", func(h *colly.HTMLElement) {
		//fmt.Println(h.ChildAttrs("div", "class"))
		//fmt.Println(h.ChildAttr("div[class=partShipping]", "class"))
		//fmt.Println(h.ChildText("div[class=partShipping]"))
		name := strings.Split(h.Response.Request.URL.String(), "/")
		price := h.ChildText("div[class=partPrice]")
		shipping := h.ChildText("div[class=partShipping]")
		img := h.ChildAttr("img", "src")
		grade := h.ChildText("div[class=gradeText]")
		p := Part{
			Name:     name[len(name)-1],
			URL:      h.Response.Request.URL.String(),
			Grade:    grade,
			Img:      img,
			Price:    price,
			Shipping: shipping,
		}
		parts = append(parts, p)

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	// https://www.hollanderparts.com/

	c.Visit("https://www.hollanderparts.com/" + url_suffix)

	c.Wait()

	fmt.Println(parts)

	return parts

}

func getEngines(url_suffix string) []Part {
	var engines []Part
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		//fmt.Println(h)
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			p := Part{
				Name: h.ChildText("a"),
				URL:  h.ChildAttr("a", "href"),
			}
			engines = append(engines, p)
			c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
		})

		//link := h.ChildAttr("a", "href")
		//if strings.Index(link, "engine") != -1 {
		//fmt.Println("Fetching Engine Parts...")
		////c.Visit(h.Request.AbsoluteURL(link))

		//}
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	// https://www.hollanderparts.com/
	c.Visit("https://www.hollanderparts.com/" + url_suffix)
	c.Wait()

	return engines

}

func getEngineLinks(w http.ResponseWriter, r *http.Request) {
	var links []Part
	// Data Structure is Category -> Part Type -> Part -> Part Fitment -> Mathing Parts
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			if h.ChildText("a") == "Engine" && h.ChildAttr("a", "href") != "" {
				p := Part{
					Name: h.ChildText("a"),
					URL:  h.ChildAttr("a", "href"),
				}
				links = append(links, p)
				c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
			}
		})

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	//c.Visit("https://www.hollanderparts.com/")
	params := mux.Vars(r)
	vin := params["vin"]
	c.PostMultipart("https://www.hollanderparts.com/Home", map[string][]byte{
		"hdnVIN": []byte(vin),
	})
	c.Wait()

	w.Header().Add("Content-Type", "application/json")

	engineLink := links[1]

	engines := getEngines(engineLink.URL)

	for _, engine := range engines {
		getParts(engine.URL)
	}

	j, err := json.Marshal(parts)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(j)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{vin}", getEngineLinks)

	fmt.Println("Running on localhost:8000/{vin}")
	http.ListenAndServe(":8000", router)
}

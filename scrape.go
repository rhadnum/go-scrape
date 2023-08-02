package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"encoding/json"
	"io/ioutil"
	"log"
)

type Page struct {
	Url string
	ScrapedUrls []string
}

func main() {
	currentLink := "https://en.wikipedia.org/wiki/Web_scraping"
	
    c := colly.NewCollector()
	pages := []Page{}
	
	c.OnRequest(func(r *colly.Request) {
		currentLink = r.URL.String()
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		content, err := ioutil.ReadFile("scrape.json")
		if err != nil {
			log.Fatal(err)
		}

		link := e.Attr("href")

		readPages := []Page{}
		err = json.Unmarshal(content, &readPages)

		if len(readPages) > 0 {
			for i := range readPages {
				if readPages[i].Url == currentLink {
					readPages[i].ScrapedUrls = append(readPages[i].ScrapedUrls, link)
					pages = readPages
				} else {
					var page = Page{ currentLink, []string{} }
					pages = append(pages, page)
				}
			}
		} else {
			var page = Page{ currentLink, []string{} }
			pages = append(pages, page)
		}
		
		data, _ := json.Marshal(pages)

		err = ioutil.WriteFile("scrape.json", data, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// TODO - Follow links to other pages
		// e.Request.Visit(link)
    })

    c.Visit(currentLink)
}
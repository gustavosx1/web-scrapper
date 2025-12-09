package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gocolly/colly/v2"
)

type LinkInfo struct {
	URL    *url.URL
	Status int
}

type LinkList []LinkInfo

func (l LinkList) Print(title string) {
	fmt.Println("------------------------------------")
	fmt.Println(title)
	fmt.Println("------------------------------------")

	for _, r := range l {
		fmt.Printf("\nLink: %v", r.URL)
		fmt.Printf("\nStatus Code: %v\n", r.Status)
	}
}

func main() {
	var goodLinks LinkList
	var badLinks LinkList

	c := colly.NewCollector(
		colly.AllowedDomains("scrape-me.dreamsofcode.io"),
	)

	// Coletar cada link encontrado
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		absURL := e.Request.AbsoluteURL(href)

		// evita visitas repetidas
		if has, _ := c.HasVisited(absURL); !has {
			c.Visit(absURL)
		}
	})

	// Respostas com status 200
	c.OnResponse(func(r *colly.Response) {
		u := r.Request.URL

		if r.StatusCode == 200 {
			goodLinks = append(goodLinks, LinkInfo{
				URL:    u,
				Status: r.StatusCode,
			})
		} else {
			badLinks = append(badLinks, LinkInfo{
				URL:    u,
				Status: r.StatusCode,
			})
		}
	})

	// Respostas com erro (404, 500, etc)
	c.OnError(func(r *colly.Response, err error) {
		badLinks = append(badLinks, LinkInfo{
			URL:    r.Request.URL,
			Status: r.StatusCode,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando:", r.URL.String())
	})

	// Inicia o scraping
	err := c.Visit("https://scrape-me.dreamsofcode.io/")
	if err != nil {
		log.Println("Erro:", err)
	}

	goodLinks.Print("Links que funcionam dentro do site:")
	badLinks.Print("Links que NÃO funcionam dentro do site:")
}

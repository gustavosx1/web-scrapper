// Package scrapper
package scrapper

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

type LinkInfo struct {
	URL    *url.URL
	Status int
}

type LinkList []LinkInfo

func Scrapper(url string) (LinkList, LinkList) {
	var goodLinks LinkList
	var badLinks LinkList

	dominios := strings.Split(url, "/")
	c := colly.NewCollector(
		colly.AllowedDomains(dominios[2]),
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
	err := c.Visit(url)
	if err != nil {
		log.Println("Erro:", err)
	}

	return goodLinks, badLinks
}

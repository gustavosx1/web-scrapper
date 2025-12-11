package main

import (
	"go-scrapper/request"
	api "go-scrapper/response"
)

func main() {
	URL := "http://localhost:7070/scrape?url=https://scrape-me.dreamsofcode.io/"

	api.API()

	request.Request(URL)
	// depois de terminado que eu parei pra pensar que isso aqui não faz sentido algum
}

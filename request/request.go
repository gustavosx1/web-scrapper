// Package request
package request

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Request(URL string) {
	Client := http.Client{
		Timeout: 25 * time.Second,
	}
	req, _ := http.NewRequest("POST", URL, nil)
	resp, err := Client.Do(req)
	if err != nil {
		http.Error(nil, err.Error(), http.StatusBadRequest)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao parsear o body", err)
		return
	}
	fmt.Println(string(body))
}

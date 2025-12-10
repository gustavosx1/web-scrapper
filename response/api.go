// Package api
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	scrapper "go-scrapper/model"
)

// Estrutura da resposta JSON
type APIResponse struct {
	GoodLinks []LinkItem `json:"good_links"`
	BadLinks  []LinkItem `json:"bad_links"`
}

type LinkItem struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

// App principal
type Application struct {
	Logger *slog.Logger
}

func API() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	app := &Application{Logger: logger}

	http.HandleFunc("/scrape", app.ScrapeHandler)

	PORT := ":7070"

	app.Logger.Info("Servidor iniciado", slog.String("port", PORT))

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Error("Erro ao iniciar servidor", slog.String("error", err.Error()))
	}
}

func (app *Application) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// pegar: /scrape?url=xxx
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "Parâmetro 'url' é obrigatório", http.StatusBadRequest)
		return
	}
	/*
		u, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, "Erro ao fazer o parsing do url RAW", http.StatusBadRequest)
		}
	*/
	app.Logger.Info("Scraping iniciado", slog.String("url", targetURL))

	// Scrape
	good, bad := scrapper.Scrapper(targetURL)

	// converter para resposta JSON
	resp := APIResponse{
		GoodLinks: convertLinks(good),
		BadLinks:  convertLinks(bad),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Converte suas structs internas para JSON amigável
func convertLinks(list []scrapper.LinkInfo) []LinkItem {
	result := make([]LinkItem, 0, len(list))
	for _, l := range list {
		result = append(result, LinkItem{
			URL:    l.URL.String(),
			Status: l.Status,
		})
	}
	return result
}

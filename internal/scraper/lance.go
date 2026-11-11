package scraper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

const SourceLance = "Lance!"

type LanceScraper struct {
	URL string
}

var _ SiteScraper = (*LanceScraper)(nil)

// Estruturas auxiliares para decodificar o JSON-LD do Lance!
type LdListItem struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

type LdMainEntity struct {
	Type              string       `json:"@type"`
	ItemListElement   []LdListItem `json:"itemListElement"`
}

type LdGraphItem struct {
	Type       string        `json:"@type"`
	MainEntity *LdMainEntity `json:"mainEntity"`
}

type LdSchema struct {
	Context string        `json:"@context"`
	Graph   []LdGraphItem `json:"@graph"`
}

func (l *LanceScraper) Fetch() ([]News, error) {
	var newsList 	[]News
	var scraperErr 	error

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ Resposta recebida de [%s] [Status: %d]\n", r.Request.URL, r.StatusCode)
	})

	collector.OnError(func(r *colly.Response, err error){
		scraperErr = fmt.Errorf("erro HTTP no scraper do Lance! [Status: %d]: %w", r.StatusCode, err)
	})

	collector.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		content := e.Text

		if strings.Contains(content, "ItemList") && strings.Contains(content, "ListItem") {
			var schema LdSchema
			err := json.Unmarshal([]byte(content), &schema)
			if err == nil {
				for _, graph := range schema.Graph {
					if graph.MainEntity != nil && graph.MainEntity.Type == "ItemList" {
						for _, item := range graph.MainEntity.ItemListElement {
							if item.Name != "" && item.URL != "" {
								if strings.Contains(item.URL, "/vasco") || strings.Contains(item.URL, "/futebol-nacional/") {
									if !strings.HasSuffix(item.URL, "/vasco") && !strings.HasSuffix(item.URL, "/vasco/") {
										newsList = append(newsList, News{
											Title:  strings.TrimSpace(item.Name),
											Link:   strings.TrimSpace(item.URL),
											Source: SourceLance,
										})
									}
								}
							}
						}
					}
				}
			}
		}
	})

	err := collector.Visit(l.URL)
	if err != nil {
		return nil, fmt.Errorf("falha ao visitar Lance!: %w", err)
	}

	if scraperErr != nil {
		return nil, scraperErr
	}

	uniqueNews := RemoveDuplicateNews(newsList)
	fmt.Printf("🔍 Lance! Scraper encontrou %d notícias únicas.\n", len(uniqueNews))
	return uniqueNews, nil
}
package scraper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

type LanceScraper struct {
	URL string
}

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
	var newsList []News

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// Captura o conteúdo da tag script que contém o ItemList estruturado do JSON-LD
	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		content := e.Text

		// Filtramos apenas o bloco JSON que possui a lista de itens (ItemList)
		if strings.Contains(content, "ItemList") && strings.Contains(content, "ListItem") {
			var schema LdSchema
			err := json.Unmarshal([]byte(content), &schema)
			if err == nil {
				for _, graph := range schema.Graph {
					if graph.MainEntity != nil && graph.MainEntity.Type == "ItemList" {
						for _, item := range graph.MainEntity.ItemListElement {
							if item.Name != "" && item.URL != "" {
								// Filtra para garantir que foca em matérias do Vasco ou do contexto do clube
								if strings.Contains(item.URL, "/vasco") || strings.Contains(item.URL, "/futebol-nacional/") {
									// Evita pegar links genéricos da home da tag
									if !strings.HasSuffix(item.URL, "/vasco") && !strings.HasSuffix(item.URL, "/vasco/") {
										newsList = append(newsList, News{
											Title:  strings.TrimSpace(item.Name),
											Link:   strings.TrimSpace(item.URL),
											Source: "Lance!",
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

	err := c.Visit(l.URL)
	if err != nil {
		return nil, err
	}

	uniqueNews := removeDuplicateNews(newsList)
	fmt.Printf("🔍 Lance! Scraper encontrou %d notícias via JSON-LD.\n", len(uniqueNews))
	return uniqueNews, nil
}
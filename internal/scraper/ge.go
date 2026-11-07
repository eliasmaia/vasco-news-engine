package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

type GeScraper struct {
	URL string
}

func (g *GeScraper) Fetch() ([]News, error) {
	var newsList []News

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// Log para sabermos se a requisição HTTP chegou bem
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ Resposta recebida do GE [%s] [Status: %d]\n", g.URL, r.StatusCode)
	})

	// Captura erros de requisição
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("❌ Erro no scraper do GE: %v (Status: %d)\n", err, r.StatusCode)
	})

	// Vamos tentar múltiplos seletores comuns do GE caso a classe principal tenha mudado
	c.OnHTML("a.feed-post-link", func(e *colly.HTMLElement) {
		titleElement := e.Text
		linkElement := e.Attr("href")

		if titleElement != "" && linkElement != "" {
			link := linkElement
			if !strings.HasPrefix(link, "http") {
				link = "https:" + link
			}

			newsList = append(newsList, News{
				Title:  strings.TrimSpace(titleElement),
				Link:   link,
				Source: "GloboEsporte",
			})
		}
	})

	err := c.Visit(g.URL)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔍 GE Scraper encontrou %d notícias brutas.\n", len(newsList))
	return newsList, nil
}
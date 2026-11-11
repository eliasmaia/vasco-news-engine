package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

const SourceSupervasco = "Supervasco"

type Supervasco struct{
		URL	string
}

var _ SiteScraper = (*Supervasco)(nil)

func (s *Supervasco) Fetch() ([]News, error) {
	var newsList 	[]News
	var scraperErr 	error

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ Resposta recebida de [%s] [Status: %d]\n", r.Request.URL, r.StatusCode)
	})

	collector.OnError(func(r *colly.Response, err error){
		scraperErr = fmt.Errorf("erro HTTP no scraper do Supervasco [Status: %d]: %w", r.StatusCode, err)
	})

	collector.OnHTML("a", func(e *colly.HTMLElement) {
		linkElement := e.Attr("href")
		titleElement := strings.TrimSpace(e.Text)

		if strings.Contains(linkElement, "/noticias/") && len(titleElement) > 20 {
			newsList = append(newsList, News{
				Title:  titleElement,
				Link:   e.Request.AbsoluteURL(linkElement),
				Source: "Supervasco",
			})
		}
	})

	err := collector.Visit(s.URL)
	if err != nil {
		return nil, fmt.Errorf("falha ao visitar Supervasco: %w", err)
	}

	if scraperErr != nil {
		return nil, scraperErr
	}

	uniqueNews := RemoveDuplicateNews(newsList)
	fmt.Printf("🔍 Supervasco Scraper encontrou %d notícias únicas.\n", len(uniqueNews))
	return uniqueNews, nil
}

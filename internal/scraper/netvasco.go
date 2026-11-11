package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

const SourceNetVasco = "NetVasco"

type NetvascoScraper struct {
	URL string
}

var _ SiteScraper = (*NetvascoScraper)(nil)

func (n *NetvascoScraper) Fetch() ([]News, error) {
	var newsList 	[]News
	var scraperErr 	error

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ Resposta recebida de [%s] [Status: %d]\n", r.Request.URL, r.StatusCode)
	})

	collector.OnError(func(r *colly.Response, err error){
		scraperErr = fmt.Errorf("erro HTTP no scraper do NetVasco [Status: %d]: %w", r.StatusCode, err)
	})

	c.OnHTML("a", func(el *colly.HTMLElement) {
		href := el.Attr("href")
		titleElement := strings.TrimSpace(el.Text)

		if href == "" || titleElement == "" {
			return
		}

		linkElement := el.Request.AbsoluteURL(href)

		if strings.Contains(linkElement, "netvasco.com.br") && len(titleElement) > 20 {
			if !strings.HasSuffix(linkElement, "netvasco.com.br/") && 
			   !strings.Contains(linkElement, "esportes.html") && 
			   !strings.Contains(linkElement, "contato") && 
			   !strings.Contains(linkElement, "expediente") {
				
				newsList = append(newsList, News{
					Title:  titleElement,
					Link:   linkElement,
					Source: SourceNetVasco,
				})
			}
		}
	})

	err := c.Visit(n.URL)
	if err != nil {
		return nil, fmt.Errorf("falha ao visitar NetVasco: %w", err)
	}

	if scraperErr != nil {
		return nil, scraperErr
	}

	uniqueNews := RemoveDuplicateNews(newsList)
	fmt.Printf("🔍 NetVasco Scraper encontrou %d notícias únicas.\n", len(uniqueNews))
	return uniqueNews, nil
}
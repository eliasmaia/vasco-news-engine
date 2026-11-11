package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

const SourceGloboEsporte = "GloboEsporte"

type GeScraper struct {
	URL string
}

var _ SiteScraper = (*GeScraper)(nil)

func (g *GeScraper) Fetch() ([]News, error) {
	var newsList 	[]News
	var scraperErr 	error

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ Resposta recebida de [%s] [Status: %d]\n", r.Request.URL, r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		scraperErr = fmt.Errorf("erro HTTP no scraper do GE [Status: %d]: %w", r.StatusCode, err)
	})

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
				Source: SourceGloboEsporte,
			})
		}
	})

	err := c.Visit(g.URL)
	if err != nil {
		return nil, fmt.Errorf("falha ao visitar GE: %w", err)
	}

	if scraperErr != nil {
		return nil, scraperErr
	}

	uniqueNews := RemoveDuplicateNews(newsList)
	fmt.Printf("🔍 GE Scraper encontrou %d notícias únicas.\n", len(uniqueNews))
	return uniqueNews, nil
}
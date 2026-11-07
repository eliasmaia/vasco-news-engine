package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

type NetvascoScraper struct {
	URL string
}

func (n *NetvascoScraper) Fetch() ([]News, error) {
	var newsList []News

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.OnHTML("a", func(el *colly.HTMLElement) {
		href := el.Attr("href")
		title := strings.TrimSpace(el.Text)

		if href == "" || title == "" {
			return
		}

		// Converte para URL absoluta primeiro para garantir que o domínio seja avaliado corretamente
		absURL := el.Request.AbsoluteURL(href)

		// Filtra links que pertencem ao NetVasco e têm tamanho de título compatível com notícia
		if strings.Contains(absURL, "netvasco.com.br") && len(title) > 20 {
			// Ignora a home pura, páginas institucionais e menus fixos
			if !strings.HasSuffix(absURL, "netvasco.com.br/") && 
			   !strings.Contains(absURL, "esportes.html") && 
			   !strings.Contains(absURL, "contato") && 
			   !strings.Contains(absURL, "expediente") {
				
				newsList = append(newsList, News{
					Title:  title,
					Link:   absURL,
					Source: "NetVasco",
				})
			}
		}
	})

	err := c.Visit(n.URL)
	if err != nil {
		return nil, err
	}

	uniqueNews := removeDuplicateNews(newsList)
	fmt.Printf("🔍 NetVasco Scraper encontrou %d notícias brutas.\n", len(uniqueNews))
	return uniqueNews, nil
}
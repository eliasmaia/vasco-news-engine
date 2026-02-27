package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Supervasco struct{}

func (supervasco *Supervasco) Fetch() ([]News, error) {
	var newsList []News
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ Resposta recebida de %s [Status: %d]\n", r.Request.URL, r.StatusCode)
	})

	collector.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		titulo := strings.TrimSpace(e.Text)

		if strings.Contains(link, "/noticias/") && len(titulo) > 20 {
			n := News{
				Title:  titulo,
				Link:   e.Request.AbsoluteURL(link),
				Source: "Supervasco",
			}
			newsList = append(newsList, n)
		}
	})

	err := collector.Visit("https://www.supervasco.com/")
	if err != nil {
		fmt.Println("❌ Erro ao visitar site:", err)
	}

	return newsList, err
}

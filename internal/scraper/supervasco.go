package scraper

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type Supervasco struct{}

func (supervasco *Supervasco) Fetch() ([]News, error) {
	var newsList []News
	collector := colly.NewCollector()

	collector.OnHTML(".noticia-titulo a", func(e *colly.HTMLElement) {
		n := News{
			Title:  strings.TrimSpace(e.Text),
			Link:   e.Request.AbsoluteURL(e.Attr("href")),
			Source: "Supervasco",
		}
		newsList = append(newsList, n)
	})

	err := collector.Visit("https://www.supervasco.com/")
	return newsList, err
}

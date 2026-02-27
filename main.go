package main

import (
	"fmt"
	"sync"
	"vasco-news-engine/internal/scraper"
)

func main() {
	scrapers := []scraper.SiteScraper{
		&scraper.Supervasco{},
	}

	var wg sync.WaitGroup
	results := make(chan scraper.News)

	fmt.Println("Iniciando busca de noticias do Vascão...")

	for _, s := range scrapers {
		wg.Add(1)
		go func(src scraper.SiteScraper) {
			defer wg.Add(-1)
			news, err := src.Fetch()
			if err == nil {
				for _, n := range news[:5] {
					results <- n
				}
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for n := range results {
		fmt.Printf("[%s] %s\n🔗 %s\n\n", n.Source, n.Title, n.Link)
	}

}

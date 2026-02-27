package main

import (
	"fmt"
	"log"
	"sync"
	"vasco-news-engine/internal/bot"
	"vasco-news-engine/internal/scraper"
	"vasco-news-engine/internal/storage"
)

func main() {

	//configuracoes, idealmente em um .env
	const (
		botToken = "TOKEN"
		chatId   = 123456
		dbPath   = "./vasco_news.db"
	)

	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatal("Erro no banco:", err)
	}

	tgBot, err := bot.NewTelegramBot(botToken, chatId)
	if err != nil {
		log.Fatal("Erro no telegram:", err)
	}

	scrapers := []scraper.SiteScraper{
		&scraper.Supervasco{},
	}

	var wg sync.WaitGroup
	results := make(chan scraper.News)

	fmt.Println("💢 Buscando notícias do Vascão...")

	for _, s := range scrapers {
		wg.Add(1)
		go func(src scraper.SiteScraper) {
			defer wg.Done()
			news, err := src.Fetch()
			if err == nil {
				for _, n := range news {
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
		if db.IsNew(n.Link) {
			fmt.Printf("🆕 Enviando: %s\n", n.Title)
			err := tgBot.SendNews(n.Title, n.Link, n.Source)
			if err == nil {
				db.Save(n.Link, n.Title, n.Source)
			} else {
				fmt.Println("❌ Erro ao enviar Telegram:", err)
			}
		}
	}

}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"vasco-news-engine/internal/bot"
	"vasco-news-engine/internal/scraper"
	"vasco-news-engine/internal/storage"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}
	//configuracoes, idealmente em um .env

	token := os.Getenv("TELEGRAM_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	dbPath := os.Getenv("DB_PATH")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Erro: TELEGRAM_CHAT_ID deve ser um número")
	}

	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatal("Erro no banco:", err)
	}

	tgBot, err := bot.NewTelegramBot(token, chatID)
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

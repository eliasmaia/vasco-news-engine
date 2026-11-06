package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"vasco-news-engine/internal/bot"
	"vasco-news-engine/internal/scraper"
	"vasco-news-engine/internal/storage"

	"github.com/joho/godotenv"
)

type App struct {
	DB       *storage.DB
	Bot      *bot.TelegramBot
	Scrapers []scraper.SiteScraper
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Erro: TELEGRAM_CHAT_ID deve ser um número")
	}

	db, err := storage.NewDB()
	if err != nil {
		log.Fatal("Erro no banco:", err)
	}

	tgBot, err := bot.NewTelegramBot(token, chatID)
	if err != nil {
		log.Fatal("Erro no telegram:", err)
	}

	app := &App{
		DB:  db,
		Bot: tgBot,
		Scrapers: []scraper.SiteScraper{
			&scraper.Supervasco{},
			&scraper.OfficialScraper{URL: "https://vasco.com.br/noticias/"},
		},
	}

	intervalo := 15 * time.Minute
	ticker := time.NewTicker(intervalo)
	defer ticker.Stop()

	fmt.Println("🤖 Bot rodando... Pressione Ctrl+C para parar.")

	app.checkAndNotify()

	for {
		select {
		case t := <-ticker.C:
			fmt.Printf("🔔 Pulso recebido em: %s\n", t.Format("15:04:05"))
			app.checkAndNotify()
		}
	}

}

func (a *App) checkAndNotify() {
	fmt.Println("💢 Buscando notícias do Vascão...")

	var wg sync.WaitGroup
	results := make(chan scraper.News)

	for _, s := range a.Scrapers {
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

	novasCount := 0
	for n := range results {
		if a.DB.IsNew(n.Link) {
			fmt.Printf("🆕 Enviando: %s\n", n.Title)
			err := a.Bot.SendNews(n.Title, n.Link, n.Source)
			if err == nil {
				a.DB.Save(n.Link, n.Title, n.Source)
				novasCount++
			} else {
				fmt.Println("❌ Erro ao enviar Telegram:", err)
			}
		}
	}

	if novasCount == 0 {
		fmt.Println("😴 Nenhuma notícia nova encontrada.")
	} else {
		fmt.Printf("✅ %d novas notícias enviadas!\n", novasCount)
	}
}

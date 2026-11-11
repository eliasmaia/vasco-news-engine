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

type App struct {
	DB       *storage.DB
	Bot      *bot.TelegramBot
	Scrapers []scraper.SiteScraper
}

func main() {
	// try to obtain the .env local file; if it does not, then (GitHub Actions)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	token 		:= os.Getenv("TELEGRAM_TOKEN")
	chatIDStr 	:= os.Getenv("TELEGRAM_CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Erro: TELEGRAM_CHAT_ID deve ser um número válido")
	}

	db, err := storage.NewDB()
	if err != nil {
		log.Fatal("Erro crítico ao conectar no banco:", err)
	}

	tgBot, err := bot.NewTelegramBot(token, chatID)
	if err != nil {
		log.Fatal("Erro crítico ao inicializar bot do telegram:", err)
	}

	app := &App{
		DB:  db,
		Bot: tgBot,
		Scrapers: []scraper.SiteScraper{
			&scraper.Supervasco{URL: "https://www.supervasco.com"},
			&scraper.GeScraper{URL: "https://ge.globo.com/futebol/times/vasco/"},
			&scraper.NetvascoScraper{URL: "https://www.netvasco.com.br/"}, // <-- NetVasco adicionado no lugar do UOL
			&scraper.LanceScraper{URL: "https://www.lance.com.br/vasco"},
		},
	}

	fmt.Println("🤖 Iniciando verificação de notícias do Vascão...")

	// runs the check once, and the program finishes naturally (exit 0)
	app.checkAndNotify()
}

func (a *App) checkAndNotify() {
	fmt.Println("💢 Buscando notícias do Vascão em paralelo...")

	var wg sync.WaitGroup
	results := make(chan scraper.News)
	
	//dispatches the scrapers concurrently using the Goroutines
	for _, s := range a.Scrapers {
		wg.Add(1)
		go func(src scraper.SiteScraper) {
			defer wg.Done()
			news, err := src.Fetch()
			if err != nil {
				fmt.Printf("⚠️ Alerta no scraper: %v\n", err)
				return
			}
			for _, n := range news {
				results <- n
			}
		}(s)
	}

	// goroutine to close all channels once all threads finish
	go func() {
		wg.Wait()
		close(results)
	}()

	newsCount := 0
	for n := range results {
		if a.DB.IsNew(n.Link) {
			fmt.Printf("🆕 Noticia inédita encontrada: %s\n", n.Title)
			
			msgDTO := bot.NewsMessage{
				Title:	n.Title,
				Link:	n.Link,
				Source:	n.Source,
			}

			err := a.Bot.SendNews(msgDTO)
			if err == nil {
				saveErr := a.DB.Save(n.Link, n.Title, n.Source)
				if saveErr != nil {
					fmt.Printf("❌ Erro ao salvar a noticia no banco: %v\n", saveErr)	
				} else {
					newsCount++
				}
			} else {
				fmt.Printf("❌ Erro ao enviar mensagem para o Telegram: %v\n", err)
			}
		}
	}

	if newsCount == 0 {
		fmt.Println("😴 Nenhuma notícia nova encontrada.")
	} else {
		fmt.Printf("✅ %d novas notícias enviadas com sucesso!\n", newsCount)
	}
}
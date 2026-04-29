# 💢 Vasco News Engine

[<img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/us.svg" width="20"> English](#english) |[ <img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/br.svg" width="20"> Português](#português)


<a name="english"></a>
## <img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/us.svg" width="25"> English

Vasco News Engine is an automated monitoring service built with Go. It periodically crawls Vasco da Gama news portals and sends real-time notifications via Telegram, ensuring you never miss an update about the club.

### 🛠️ Tech Stack
Go 1.25: High-performance concurrency using Goroutines, Channels, and WaitGroups.

SQLite (Modernc): CGO-free data persistence to prevent duplicate notifications.

Colly & GoQuery: Industry-standard web scraping and HTML parsing libraries.

Telebot v3: Robust framework for Telegram Bot API integration.

### 🏗️ Architecture
The project follows a modular design for scalability:

/bot: Handles Telegram message delivery and Markdown formatting.

/scraper: Unified interface allowing easy addition of new news sources.

/storage: Persistence layer to track previously sent links.

### 🚀 Getting Started
Create a .env file:

Snippet de código

TELEGRAM_TOKEN=your_token
TELEGRAM_CHAT_ID=your_id
DB_PATH=vasco.db
Run the engine:

Bash

go run main.go

<a name="português"></a>
## <img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/br.svg" width="25"> Português

O Vasco News Engine é um serviço de monitoramento automatizado desenvolvido em Go. Ele realiza buscas periódicas em portais de notícias do Vasco da Gama e notifica novos conteúdos via Telegram, garantindo que você nunca perca uma atualização do Gigante da Colina.

### 🛠️ Tecnologias
Go 1.25: Utilizando Goroutines, Channels e WaitGroups para scraping concorrente.

SQLite (Modernc): Persistência de dados CGO-free para evitar notificações duplicadas.

Colly & GoQuery: Ferramentas robustas para extração de dados web.

Telebot v3: Framework para integração com a API do Telegram.

### 🏗️ Arquitetura
O projeto utiliza uma estrutura modular para facilitar a manutenção:

/bot: Gerencia o envio de mensagens e formatação em Markdown.

/scraper: Interface unificada que permite adicionar novas fontes de notícias facilmente.

/storage: Camada de persistência para controle de histórico.

### 🚀 Como Rodar
Configure seu arquivo .env:

Snippet de código

TELEGRAM_TOKEN=seu_token
TELEGRAM_CHAT_ID=seu_id
DB_PATH=vasco.db
Comando para executar:

Bash

go run main.go


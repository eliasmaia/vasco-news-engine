# 💢 Vasco News Engine

[<img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/us.svg" width="20"> English](#english) |[ <img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/br.svg" width="20"> Português](#português)


<a name="english"></a>
## <img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/us.svg" width="25"> English

**Vasco News Engine** is an automated monitoring service built with Go. It concurrently crawls Vasco da Gama news portals (GloboEsporte, Lance!, NetVasco and SuperVasco) and sends real-time notifications via Telegram, ensuring you never miss an update about the club.

### 🛠️ Tech Stack
- **Go 1.25**: High-performance concurrency using Goroutines, Channels, and WaitGroups.
- **Turso DB (`libsql`)**: Cloud SQLite data persistence to prevent duplicate notifications.
- **Colly & `encoding/json`**: Web scraping and structured JSON metadata parsing.
- **Telebot v3**: Robust framework for Telegram Bot API integration.
- **GitHub Actions & Cron-job.org**: Headless cloud execution and automated scheduling.

### 🏗️ Architecture & Automation
The project follows a modular design for scalability and maintainability:

- `/bot`: Handles Telegram message delivery, DTO mapping and Markdown formatting.
- `/scraper`: Unified interface (`SiteScraper`) allowing easy addition of new news sources.
- `/storage`: Persistence layer to track previously sent links.
- **`/.github/workflows/bot.yml`**: GitHub Actions workflow triggered via `repository_dispatch` or `workflow_dispatch`, configured with caching for fast Go builds and secure environment secrets injection. External cron triggers (such as `cron-job.org`) can invoke this workflow periodically to keep the engine running seamlessly.

### 🚀 Getting Started
1. **Create a .env file** in the root directory:
```Code Snippet
TELEGRAM_TOKEN=your_token
TELEGRAM_CHAT_ID=your_id
TURSO_DATABASE_URL=libsql://your-database.turso.io
TURSO_AUTH_TOKEN=your_token
```
2. **Run the engine**:

```Bash
go run main.go
```

<a name="português"></a>
## <img src="https://raw.githubusercontent.com/lipis/flag-icons/main/flags/4x3/br.svg" width="25"> Português

O **Vasco News Engine** é um serviço de monitoramento automatizado desenvolvido em Go. Ele realiza buscas concorrentes em portais de notícias do Vasco da Gama (GloboEsporte, Lance!, NetVasco e SuperVasco) e notifica novos conteúdos via Telegram, garantindo que você nunca perca uma atualização do Gigante da Colina.

### 🛠️ Tecnologias
- **Go 1.25**: Utilizando Goroutines, Channels e WaitGroups para scraping concorrente de alta performance.
- **Turso DB (`libsql`)**: Persistência de dados em nuvem para evitar notificações duplicadas.
- **Colly & `encoding/json`**: Extração de dados web e decodificação de metadados JSON estruturados.
- **Telebot v3**: Framework moderno para integração com a API do Telegram.
- **GitHub Actions & Cron-job.org**: Execução em nuvem sem servidor dedicado e agendamento automatizado.

### 🏗️ Arquitetura
O projeto utiliza uma estrutura modular para facilitar a manutenção e escalabilidade:

- `/bot`: Gerencia o envio de mensagens, mapeamento de DTOs e formatação em Markdown.
- `/scraper`: Interface unificada (SiteScraper) que permite adicionar novas fontes de notícias facilmente.
- `/storage`: Camada de persistência para controle de histórico de links enviados.
- **`/.github/workflows/bot.yml`**: Workflow do GitHub Actions configurado para aceitar disparos via `repository_dispatch` ou `workflow_dispatch`. Ele gerencia cache de módulos Go, injeta as credenciais de forma segura via GitHub Secrets e executa o bot em ambiente `ubuntu-latest`. Serviços externos como o `cron-job.org` podem ser configurados para chamar o webhook do GitHub periodicamente, mantendo o bot rodando de forma 100% automatizada na nuvem.

### 🚀 Como Rodar
1. **Configure seu arquivo .env** na raiz do projeto:

```Snippet de código
TELEGRAM_TOKEN=seu_token
TELEGRAM_CHAT_ID=seu_id
TURSO_DATABASE_URL=libsql://seu-banco.turso.io
TURSO_AUTH_TOKEN=seu_token
```

2. **Comando para executar**:

```Bash
go run main.go
```
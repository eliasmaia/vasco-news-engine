package bot

import (
	"time"

	"gopkg.in/telebot.v3"
)

//BotClient defines the contract (interface) for any bot client (ISP applied)
type BotClient interface {
	Send(to telebot.Recipient, what interface{}, options ...interface{}) (*telebot.Message, error)
}

//TelegramBot manages the sending infrastructure (DIP applied)
type TelegramBot struct {
	client BotClient
	ChatID int64
}

//NewTelegramBot works like the Factory Method, to initialize the bot
func NewTelegramBot(token string, chatID int64) (*TelegramBot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{client: b, ChatID: chatID}, nil
}

//SendNews receives the message DTO and send it to the client injected
func (t *TelegramBot) SendNews(news NewsMessage) error {
	msg := news.Format()
	_, err := t.client.Send(telebot.ChatID(t.ChatID), msg, telebot.ModeMarkdown)
	return err
}

package bot

import (
	"time"

	"gopkg.in/telebot.v3"
)

type TelegramBot struct {
	Bot    *telebot.Bot
	ChatID int64
}

func NewTelegramBot(token string, chatID int64) (*TelegramBot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{Bot: b, ChatID: chatID}, nil
}

func (t *TelegramBot) SendNews(title, link, source string) error {
	msg := "💢 *NOVIDADE NO GIGANTE* 💢\n\n" +
		title + "\n\n" +
		"📍 Fonte: " + source + "\n" +
		"🔗 [Leia mais](" + link + ")"

	_, err := t.Bot.Send(telebot.ChatID(t.ChatID), msg, telebot.ModeMarkdown)
	return err
}

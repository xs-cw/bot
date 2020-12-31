package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"net/url"
)

func NewTelegramBot(apiToken string) {
	// if need the ladder
	uRL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyURL(uRL),
	}}
	bot, err := tgbotapi.NewBotAPIWithClient(apiToken, client)
	if err != nil {
		log.Println(err)
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("get a message: [%s]: %s", update.Message.From.String(), update.Message.Text)
		// reply the message like a repeater
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Bot    *tgbotapi.BotAPI
	Config *Config
	Done   chan struct{}
}

func NewTelegram(config *Config) *Telegram {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil
	}

	bot.Debug = config.Debug

	return &Telegram{
		Bot:    bot,
		Config: config,
		Done:   make(chan struct{}),
	}
}

func (t *Telegram) StopBot() {
	t.Bot.StopReceivingUpdates()
	<-t.Done
}

func (t *Telegram) StartBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.Bot.GetUpdatesChan(u)

	go func(updates tgbotapi.UpdatesChannel) {
		for update := range updates {
			if update.Message != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				_, err := t.Bot.Send(msg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		t.Done <- struct{}{}
	}(updates)
}

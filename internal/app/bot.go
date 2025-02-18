package app

import (
	"att-diplom/internal/appinit"
	"att-diplom/internal/types"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotWrapper struct {
	*types.Bot
}

func NewApp(ctx context.Context) (*BotWrapper, error) {
	a := &BotWrapper{Bot: &types.Bot{}}

	err := appinit.InitDeps(ctx, a.Bot)
	if err != nil {
		return nil, err
	}

	return a, nil
}

var BotTelegram *tgbotapi.BotAPI

func (a *BotWrapper) RunBot() error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("необходимо установить TELEGRAM_BOT_TOKEN")
	}

	bot, err := initBot(botToken)

	if err != nil {
		return fmt.Errorf("bot initialization error: %v", err)
	}

	BotTelegram = bot

	return nil
}

func initBot(botToken string) (*tgbotapi.BotAPI, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	bot, err := tgbotapi.NewBotAPIWithClient(botToken, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании бота: %v", err)
	}

	bot.Debug = true
	return bot, nil
}

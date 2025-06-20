package app

import (
	"att-diplom/internal/handlers"
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotWrapper struct {
	EngineerBot *tgbotapi.BotAPI
	TeacherBot  *tgbotapi.BotAPI
}

func NewApp(ctx context.Context) (*BotWrapper, error) {
	botToken1 := os.Getenv("TELEGRAM_BOT_TOKEN_ENGINEER")
	botToken2 := os.Getenv("TELEGRAM_BOT_TOKEN_TEACHER")

	if botToken1 == "" || botToken2 == "" {
		return nil, fmt.Errorf("необходимо установить TELEGRAM_BOT_TOKEN_ENGINEER и TELEGRAM_BOT_TOKEN_TEACHER")
	}

	bot1, err := initBot(botToken1)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации первого бота: %v", err)
	}

	bot2, err := initBot(botToken2)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации второго бота: %v", err)
	}

	a := &BotWrapper{
		EngineerBot: bot1,
		TeacherBot:  bot2,
	}

	err = a.SetTeacherBotCommands()
	if err != nil {
		log.Printf("Ошибка установки команд для второго бота: %v", err)
	}

	return a, nil
}

var BotTelegramEngineer *tgbotapi.BotAPI

func (a *BotWrapper) RunBots(db *sql.DB) error {
	go a.listenUpdates(a.EngineerBot, "Первый бот", db)

	go a.listenUpdates(a.TeacherBot, "Второй бот", db)

	select {}
}

func (a *BotWrapper) listenUpdates(bot *tgbotapi.BotAPI, botName string, db *sql.DB) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	fmt.Printf("[%s] Бот запущен\n", botName)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handlers.HandleCommand(bot, update.Message)
			}
		}
		if update.CallbackQuery != nil {
			handlers.HandleClick(bot, update, db)
		}
		if update.Message != nil && !update.Message.From.IsBot && !update.Message.IsCommand() {
			handlers.HandleUserMessage(update.Message)
		}

	}
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

func (a *BotWrapper) SetTeacherBotCommands() error {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запуск бота"},
		{Command: "help", Description: "Справка по командам"},
		{Command: "info", Description: "Информация о боте"},
	}

	cmdCfg := tgbotapi.NewSetMyCommands(commands...)
	_, err := a.TeacherBot.Request(cmdCfg)
	if err != nil {
		return fmt.Errorf("ошибка установки команд для второго бота: %v", err)
	}
	return nil
}

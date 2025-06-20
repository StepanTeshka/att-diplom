package handlers

import (
	screenbot "att-diplom/internal/ScreenBot"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()

	log.Printf("[Учитель-бот] Получена команда: %s", command)

	var responseText string

	switch command {
	case "start":
		screenbot.MainScreenBot(bot, message)
	case "help":
		responseText = "Список команд:\n/start - Начало работы\n/help - Справка\n/info - Информация"
	case "info":
		responseText = "Этот бот создан для работы с учителями!"
	default:
		responseText = "Неизвестная команда. Напишите /help для списка команд."
	}

	msg := tgbotapi.NewMessage(chatID, responseText)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}

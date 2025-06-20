package handlers

import (
	screenbot "att-diplom/internal/ScreenBot"
	"att-diplom/internal/helper"
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleClick(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	chatID := update.CallbackQuery.Message.Chat.ID
	data := update.CallbackQuery.Data

	switch data {
	case "createApp":
		go screenbot.CreateAppScreen(bot, update.CallbackQuery.Message, db)
	case "checkStatus":
		go screenbot.CheckStatusAppScreen(bot, update.CallbackQuery.Message, db)
	case "exitCreateApp":
		helper.CancelWaitingForUser(chatID)

		bot.Send(tgbotapi.NewMessage(chatID, "❌ Создание заявки отменено."))

		screenbot.MainScreenBot(bot, update.CallbackQuery.Message)
	default:
		msg := tgbotapi.NewMessage(chatID, "Неизвестная команда.")
		bot.Send(msg)
	}

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	bot.Request(callback)
}

package screenbot

import (
	"att-diplom/internal/buttons"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainScreenBot(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Привет, я бот вычислительного центра АТТ. Чем могу помочь?")
	msg.ReplyMarkup = buttons.MainButtonsStart()
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

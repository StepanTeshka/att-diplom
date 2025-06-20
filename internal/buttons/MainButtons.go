package buttons

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func MainButtonsStart() tgbotapi.InlineKeyboardMarkup {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать заявку", "createApp"),
			tgbotapi.NewInlineKeyboardButtonData("Проверить статус заявки", "checkStatus"),
		),
	)
	return buttons
}

func ExitCreateApp() tgbotapi.InlineKeyboardMarkup {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На главную", "exitCreateApp"),
		),
	)
	return buttons
}

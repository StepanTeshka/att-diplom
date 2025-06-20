package screenbot

import (
	"att-diplom/internal/buttons"
	"att-diplom/internal/functions"
	"att-diplom/internal/helper"
	"att-diplom/internal/types"
	"database/sql"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	userProblem string
	cabinet     string
	nameTeacher string
)

func CreateAppScreen(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	chatID := message.Chat.ID

	bot.Send(tgbotapi.NewMessage(chatID, "Вы выбрали создать заявку."))

	//Название

	msg1 := tgbotapi.NewMessage(chatID, "📝 Пожалуйста, опишите вашу проблему:")
	msg1.ReplyMarkup = buttons.ExitCreateApp()
	bot.Send(msg1)

	userMsg1, err := helper.WaitForUserInput(chatID, 60*time.Second)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		return
	}

	userProblem = userMsg1.Text

	// кабинет

	msg2 := tgbotapi.NewMessage(chatID, "📍 Укажите номер кабинета:")
	msg2.ReplyMarkup = buttons.ExitCreateApp()
	bot.Send(msg2)

	userMsg2, err := helper.WaitForUserInput(chatID, 60*time.Second)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		return
	}

	cabinet = userMsg2.Text

	// инициалы

	msg3 := tgbotapi.NewMessage(chatID, "👤 Укажите свои инициалы:")
	msg3.ReplyMarkup = buttons.ExitCreateApp()
	bot.Send(msg3)

	userMsg3, err := helper.WaitForUserInput(chatID, 60*time.Second)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		return
	}

	nameTeacher = userMsg3.Text

	// запись в бд

	newApplication := types.Application{
		Description: sql.NullString{String: userProblem, Valid: true},
		NameTeacher: sql.NullString{String: nameTeacher, Valid: true},
		IDEngineer:  sql.NullInt64{Valid: false},
		Status:      sql.NullString{String: "Не назначено", Valid: true},
		StartDate:   time.Now(),
		EndDate:     sql.NullTime{Valid: false},
		Cabinet:     sql.NullString{String: cabinet, Valid: true},
	}

	newID, err := functions.AddApplication(db, newApplication)
	if err != nil {
		log.Println("Ошибка при добавление заявки: ", newID, "Причина: ", err)
		helper.CancelWaitingForUser(chatID)

		bot.Send(tgbotapi.NewMessage(chatID, "❌ Создание заявки отменено. Не смогли вас записать"))
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, "✅ Вы успешно оставили заявку, ожидайте в ближайшее время к вам подойдет инженер"))
	MainScreenBot(bot, message)

}

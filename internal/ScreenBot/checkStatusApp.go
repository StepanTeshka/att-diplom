package screenbot

import (
	"att-diplom/internal/buttons"
	"att-diplom/internal/helper"
	"att-diplom/internal/types"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckStatusAppScreen(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	chatID := message.Chat.ID

	bot.Send(tgbotapi.NewMessage(chatID, "Вы выбрали просмотреть статус заявки"))

	msg1 := tgbotapi.NewMessage(chatID, "📝 Введите номер кабинета для проверки статуса:")
	msg1.ReplyMarkup = buttons.ExitCreateApp()
	bot.Send(msg1)

	userMsg1, err := helper.WaitForUserInput(chatID, 60*time.Second)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		return
	}

	cabinet := strings.TrimSpace(userMsg1.Text)

	// поиск в бд

	var applications []types.Application

	query := `
	SELECT idapplication, description, nameteacher, idengineer, status, startdate, enddate, cabinet
	FROM applications
	WHERE cabinet = $1
`

	rows, err := db.Query(query, cabinet)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		bot.Send(tgbotapi.NewMessage(chatID, "🚫 Ошибка при получении данных."))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var app types.Application
		err := rows.Scan(
			&app.ID,
			&app.Description,
			&app.NameTeacher,
			&app.IDEngineer,
			&app.Status,
			&app.StartDate,
			&app.EndDate,
			&app.Cabinet,
		)
		if err != nil {
			log.Println("Ошибка при чтении строки:", err)
			continue
		}
		applications = append(applications, app)
	}

	if len(applications) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Заявки для этого кабинета не найдены."))
		return
	}

	var result string
	for _, app := range applications {
		entry := fmt.Sprintf(
			"📄 Заявка №%d\n👨‍🏫 Преподаватель: %s\n📦 Проблема: %s\n📍 Кабинет: %s\n📌 Статус: %s\n\n",
			app.ID,
			app.NameTeacher.String,
			app.Description.String,
			app.Cabinet.String,
			app.Status.String,
		)
		result += entry
	}

	bot.Send(tgbotapi.NewMessage(chatID, result))
	MainScreenBot(bot, message)
}

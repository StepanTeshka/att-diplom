package handlers

import (
	"att-diplom/internal/functions"
	"att-diplom/internal/helper"
	"att-diplom/internal/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, bot *tgbotapi.BotAPI) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	var application types.ApplicationUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(application.ID)
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return
	}
	startDate, err := time.Parse("2006-01-02", application.StartDate)
	if err != nil {
		http.Error(w, "Неверный формат startDate. Используйте YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var endDate sql.NullTime
	if application.EndDate != nil && *application.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", *application.EndDate)
		if err != nil {
			http.Error(w, "Неверный формат endDate. Используйте YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		endDate = sql.NullTime{Time: parsedEndDate, Valid: true}
	} else {
		endDate = sql.NullTime{Valid: false}
	}
	newApplication := types.Application{
		Description: helper.NewNullString(application.Description.String),
		NameTeacher: helper.NewNullString(application.NameTeacher.String),
		IDEngineer:  helper.NewNullInt64(application.IDEngineer.String),
		Status:      helper.NewNullString(application.Status.String),
		StartDate:   startDate,
		EndDate:     endDate,
		Cabinet:     helper.NewNullString(application.Cabinet.String),
		ID:          id,
	}
	err = functions.UpdateApplication(db, newApplication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chatIdStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIdStr == "" {
		log.Println("Ошибка: TELEGRAM_CHAT_ID не установлен")
	} else {
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			log.Println("Ошибка преобразования chatId:", err)
		} else {
			engineerText := "Не назначен"
			if application.IDEngineer.Valid && application.IDEngineer.String != "" {
				engineerID, err := strconv.Atoi(application.IDEngineer.String)
				if err != nil {
					log.Println("Ошибка преобразования ID инженера:", err)
					http.Error(w, "Некорректный ID инженера", http.StatusBadRequest)
					return
				}

				engineer, err := functions.GetEngineerByID(db, engineerID)
				if err == nil {
					engineerText = engineer.Name
				} else {
					log.Println("Инженер не найден")
				}
			}

			message := fmt.Sprintf("🚀 Обновлена заявка №%d\nОписание: %s\nПреподаватель: %s\nКабинет: %s\nСтатус: %s\nИнженер: %s",
				id, application.Description.String, application.NameTeacher.String, application.Cabinet.String, application.Status.String, engineerText)

			if err := functions.SendMessage(bot, chatId, message); err != nil {
				log.Println("Ошибка отправки сообщения в Telegram:", err)
			} else {
				log.Println("✅ Сообщение успешно отправлено в Telegram")
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

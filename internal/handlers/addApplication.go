package handlers

import (
	"att-diplom/internal/functions"
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

func AddApplicationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, bot *tgbotapi.BotAPI) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var application types.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
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
	var idEngineer sql.NullInt64
	if application.IDEngineer != nil && *application.IDEngineer != "" {
		parsedID, err := strconv.ParseInt(*application.IDEngineer, 10, 64)
		if err != nil {
			fmt.Println("Ошибка конвертации:", err)
			idEngineer = sql.NullInt64{Valid: false}
		} else {
			idEngineer = sql.NullInt64{Int64: parsedID, Valid: true}
		}
	} else {
		idEngineer = sql.NullInt64{Valid: false}
	}
	newApplication := types.Application{
		Description: sql.NullString{String: application.Description, Valid: true},
		NameTeacher: sql.NullString{String: application.NameTeacher, Valid: true},
		IDEngineer:  sql.NullInt64{Int64: idEngineer.Int64, Valid: true},
		Status:      sql.NullString{String: application.Status, Valid: true},
		StartDate:   startDate,
		EndDate:     endDate,
		Cabinet:     sql.NullString{String: application.Cabinet, Valid: true},
	}

	newID, err := functions.AddApplication(db, newApplication)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при добавлении заявки: %v", err), http.StatusInternalServerError)
		return
	}
	// телеграм уведомление
	chatIdStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIdStr == "" {
		log.Println("Ошибка: TELEGRAM_CHAT_ID не установлен")
	} else {
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			log.Println("Ошибка преобразования chatId:", err)
		} else {
			engineerText := "Не назначен"
			if application.IDEngineer != nil && *application.IDEngineer != "" {
				engineerID, err := strconv.Atoi(*application.IDEngineer)
				if err != nil {
					log.Println("Ошибка преобразования ID инженера:", err)
					http.Error(w, "Некорректный ID инженера", http.StatusBadRequest)
					return
				}

				engineer, err := functions.GetEngineerByID(db, engineerID)
				if err != nil {
					log.Println("Инженер не найден")
				}
				engineerText = engineer.Name

			}

			message := fmt.Sprintf("🚀 Новая заявка №%d\nОписание: %s\nПреподаватель: %s\nКабинет: %s\nСтатус: %s\nИнженер: %s",
				newID, application.Description, application.NameTeacher, application.Cabinet, application.Status, engineerText)

			if err := functions.SendMessage(bot, chatId, message); err != nil {
				log.Println("Ошибка отправки сообщения в Telegram:", err)
			} else {
				log.Println("✅ Сообщение успешно отправлено в Telegram")
			}
		}
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Задача добавлена с ID: %d", newID)
}

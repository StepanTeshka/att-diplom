package types

import (
	"database/sql"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

type Engineer struct {
	ID         int
	Name       string
	Email      string
	TelegramID string
}

type EngineerRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	TelegramID string `json:"telegramId"`
}

type Application struct {
	ID                int
	Description       sql.NullString
	NameTeacher       sql.NullString
	TeacherTelegramID sql.NullString
	IDEngineer        sql.NullInt64
	Status            sql.NullString
	StartDate         time.Time
	EndDate           sql.NullTime
	Cabinet           sql.NullString
}
type ApplicationUpdateRequest struct {
	ID          string         `json:"id"`
	Description sql.NullString `json:"description"`
	NameTeacher sql.NullString `json:"teacher"`
	IDEngineer  sql.NullString `json:"engineer"`
	Status      sql.NullString `json:"status"`
	StartDate   string         `json:"startDate"`
	EndDate     *string        `json:"endDate"`
	Cabinet     sql.NullString `json:"cabinet"`
}

type ApplicationRequest struct {
	Description string  `json:"description"`
	NameTeacher string  `json:"teacher"`
	IDEngineer  *string `json:"engineer"`
	Status      string  `json:"status"`
	StartDate   string  `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Cabinet     string  `json:"cabinet"`
}

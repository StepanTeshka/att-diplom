package handlers

import (
	"att-diplom/internal/functions"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendBotHandler(w http.ResponseWriter, r *http.Request, bot *tgbotapi.BotAPI) {
	chatId := int64(966872832)
	// allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	// if allowedOrigins == "" {
	// 	log.Println("необходимо установить ALLOWED_ORIGINS")
	// 	return
	// }

	// if origin := r.Header.Get("Origin"); origin != "" {
	// 	if slices.Contains(strings.Split(allowedOrigins, ","), origin) {
	// 		enableCors(&w, origin)
	// 	}
	// }

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	message := "Новый заказ"

	// Отправка сообщения через Telegram бота
	if err := functions.SendMessage(bot, chatId, message); err != nil {
		http.Error(w, "Ошибка отправки сообщения", http.StatusInternalServerError)
		return
	}

	// Ответ на успешный запрос
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сообщение успешно отправлено"))
}

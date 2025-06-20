package handlers

import (
	"att-diplom/internal/helper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUserMessage(msg *tgbotapi.Message) bool {
	if msg == nil || msg.From == nil || msg.From.IsBot {
		return false
	}

	helper.Mu.Lock()
	ch, exists := helper.WaitingUsers[msg.From.ID]
	helper.Mu.Unlock()

	if exists {
		select {
		case ch <- msg:
			return true
		default:
			return false
		}
	}
	return false
}

package helper

import (
	"errors"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	WaitingUsers = make(map[int64]chan *tgbotapi.Message)
	Mu           sync.Mutex
)

func WaitForUserInput(userID int64, timeout time.Duration) (*tgbotapi.Message, error) {
	Mu.Lock()
	ch := make(chan *tgbotapi.Message, 1)
	WaitingUsers[userID] = ch
	Mu.Unlock()

	defer func() {
		Mu.Lock()
		delete(WaitingUsers, userID)
		Mu.Unlock()
	}()

	select {
	case msg, ok := <-ch:
		if !ok {
			return nil, errors.New("❌ Ввод отменён")
		}
		return msg, nil
	case <-time.After(timeout):
		return nil, errors.New("⏰ Время ожидания истекло")
	}
}

func CancelWaitingForUser(userID int64) {
	Mu.Lock()
	defer Mu.Unlock()
	if ch, ok := WaitingUsers[userID]; ok {
		close(ch)
		delete(WaitingUsers, userID)
	}
}

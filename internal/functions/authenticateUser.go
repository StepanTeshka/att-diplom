package functions

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(db *sql.DB, email, password string) (int, error) {
	var passwordHash string
	var userID int

	query := "SELECT id, password_hash FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("пользователь не найден")
		}
		return 0, fmt.Errorf("ошибка при запросе к базе данных: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return 0, errors.New("неверный пароль")
	}

	return userID, nil
}

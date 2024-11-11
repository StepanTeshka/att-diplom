package functions

import (
	"att-diplom/internal/types"
	"database/sql"
	"fmt"
)

func GetAllEngineers(db *sql.DB) ([]types.Engineer, error) {
	var engineers []types.Engineer

	query := "SELECT idengineer, name, email, telegramid FROM engineers"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к базе данных: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var engineer types.Engineer
		err := rows.Scan(&engineer.ID, &engineer.Name, &engineer.Email, &engineer.TelegramID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных строки: %v", err)
		}
		engineers = append(engineers, engineer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return engineers, nil
}
func GetEngineerByID(db *sql.DB, id int) (types.Engineer, error) {
	var engineer types.Engineer

	query := `
		SELECT idengineer, name, email, telegramid FROM engineers WHERE idengineer = $1
	`
	err := db.QueryRow(query, id).Scan(
		&engineer.ID, &engineer.Name, &engineer.Email, &engineer.TelegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return engineer, fmt.Errorf("задача с ID %d не найдена", id)
		}
		return engineer, fmt.Errorf("ошибка при запросе к базе данных: %v", err)
	}

	return engineer, nil
}

func AddEngineer(db *sql.DB, engineer types.EngineerRequest) (int, error) {
	var newID int

	query := `
		INSERT INTO engineers (name, email, telegramId)
		VALUES ($1, $2, $3)
		RETURNING idEngineer
	`
	err := db.QueryRow(query, engineer.Name, engineer.Email, engineer.TelegramID).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении инженера в базу данных: %v", err)
	}

	return newID, nil
}

func DeleteEngineer(db *sql.DB, id int) error {
	query := "DELETE FROM engineers WHERE idEngineer = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении инженера с ID %d: %v", id, err)
	}
	return nil
}

func UpdateEngineer(db *sql.DB, engineer types.Engineer) error {
	query := `
		UPDATE engineers
		SET name = $1, email = $2, telegramId = $3
		WHERE idEngineer = $4
	`
	_, err := db.Exec(query, engineer.Name, engineer.Email, engineer.TelegramID, engineer.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных инженера с ID %d: %v", engineer.ID, err)
	}
	return nil
}

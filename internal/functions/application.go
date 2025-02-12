package functions

import (
	"att-diplom/internal/types"
	"database/sql"
	"fmt"
)

func GetApplicationsWithPagination(db *sql.DB, page, pageSize int, nameTeacher string, engineerID int, orderDate string, status string) ([]types.Application, error) {
	var applications []types.Application

	var args []interface{}
	argIndex := 1

	query := `
		SELECT idapplication, description, nameteacher, teachertelegramid, idengineer, status, startdate, enddate, cabinet
		FROM applications
		WHERE 1=1
	`

	// Добавление фильтров
	if nameTeacher != "" {
		query += fmt.Sprintf(" AND nameteacher ILIKE $%d", argIndex)
		args = append(args, "%"+nameTeacher+"%")
		argIndex++
	}

	if engineerID > 0 {
		query += fmt.Sprintf(" AND idengineer = $%d", argIndex)
		args = append(args, engineerID)
		argIndex++
	}

	if orderDate != "" {
		query += fmt.Sprintf(" AND startdate = $%d", argIndex)
		args = append(args, orderDate)
		argIndex++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Добавление параметров для пагинации
	query += fmt.Sprintf(" ORDER BY startdate LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize)
	args = append(args, (page-1)*pageSize)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к базе данных: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var application types.Application
		err := rows.Scan(&application.ID, &application.Description, &application.NameTeacher, &application.TeacherTelegramID, &application.IDEngineer, &application.Status, &application.StartDate, &application.EndDate, &application.Cabinet)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных строки: %v", err)
		}
		applications = append(applications, application)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return applications, nil
}

func GetApplicationByID(db *sql.DB, id int) (types.Application, error) {
	var application types.Application

	query := `
		SELECT idapplication, description, nameteacher, teachertelegramid, idengineer, status, startdate, enddate, cabinet
		FROM applications
		WHERE idapplication = $1
	`
	err := db.QueryRow(query, id).Scan(
		&application.ID, &application.Description, &application.NameTeacher, &application.TeacherTelegramID, &application.IDEngineer, &application.Status, &application.StartDate, &application.EndDate, &application.Cabinet)
	if err != nil {
		if err == sql.ErrNoRows {
			return application, fmt.Errorf("задача с ID %d не найдена", id)
		}
		return application, fmt.Errorf("ошибка при запросе к базе данных: %v", err)
	}

	return application, nil
}

func AddApplication(db *sql.DB, application types.Application) (int, error) {
	var newID int

	query := `
		INSERT INTO applications (description, nameteacher, idengineer, status, startdate, enddate, cabinet)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING idapplication
	`

	err := db.QueryRow(query,
		application.Description,
		application.NameTeacher,
		application.IDEngineer,
		application.Status,
		application.StartDate,
		application.EndDate,
		application.Cabinet,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи в базу данных: %v", err)
	}

	return newID, nil
}

func DeleteApplication(db *sql.DB, id int) error {
	query := "DELETE FROM applications WHERE idapplication = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи с ID %d: %v", id, err)
	}
	return nil
}

func UpdateApplication(db *sql.DB, application types.Application) error {
	query := `
		UPDATE applications
		SET description = $1, nameteacher = $2, teachertelegramid = $3, idengineer = $4, status = $5, startdate = $6, enddate = $7, cabinet = $8
		WHERE idapplication = $9
	`
	_, err := db.Exec(query,
		application.Description,
		application.NameTeacher,
		application.TeacherTelegramID,
		application.IDEngineer,
		application.Status,
		application.StartDate,
		application.EndDate,
		application.Cabinet,
		application.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи с ID %d: %v", application.ID, err)
	}
	return nil
}

package functions

import (
	"att-diplom/internal/types"
	"database/sql"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GetApplicationsWithPagination(db *sql.DB, page, pageSize int, nameTeacher string, engineerID int, orderDate string, status string) ([]types.Application, int, error) {
	var applications []types.Application
	var totalRecords int
	var args []interface{}
	argIndex := 1

	// Основной SQL-запрос с фильтрацией
	query := `
		SELECT idapplication, description, nameteacher, idengineer, status, startdate, enddate, cabinet
		FROM applications
		WHERE 1=1
	`

	countQuery := `SELECT COUNT(*) FROM applications WHERE 1=1`

	// Добавление фильтров
	if nameTeacher != "" {
		query += fmt.Sprintf(" AND nameteacher ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND nameteacher ILIKE $%d", argIndex)
		args = append(args, "%"+nameTeacher+"%")
		argIndex++
	}

	if engineerID > 0 {
		query += fmt.Sprintf(" AND idengineer = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND idengineer = $%d", argIndex)
		args = append(args, engineerID)
		argIndex++
	}

	if orderDate != "" {
		query += fmt.Sprintf(" AND startdate = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND startdate = $%d", argIndex)
		args = append(args, orderDate)
		argIndex++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status ILIKE $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Подсчет количества записей
	err := db.QueryRow(countQuery, args...).Scan(&totalRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете записей: %v", err)
	}

	// Добавление сортировки и пагинации
	query += fmt.Sprintf(" ORDER BY startdate DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при запросе к базе данных: %v", err)
	}
	defer rows.Close()

	// Читаем записи
	for rows.Next() {
		var application types.Application
		err := rows.Scan(
			&application.ID, &application.Description, &application.NameTeacher,
			&application.IDEngineer, &application.Status,
			&application.StartDate, &application.EndDate, &application.Cabinet,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при чтении данных строки: %v", err)
		}
		applications = append(applications, application)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	// Рассчитываем количество страниц
	totalPages := (totalRecords + pageSize - 1) / pageSize

	return applications, totalPages, nil
}

func GetApplicationByID(db *sql.DB, id int) (types.Application, error) {
	var application types.Application

	query := `
		SELECT idapplication, description, nameteacher, idengineer, status, startdate, enddate, cabinet
		FROM applications
		WHERE idapplication = $1
	`
	err := db.QueryRow(query, id).Scan(
		&application.ID, &application.Description, &application.NameTeacher, &application.IDEngineer, &application.Status, &application.StartDate, &application.EndDate, &application.Cabinet)
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
		SET description = $1, nameteacher = $2, idengineer = $3, status = $4, startdate = $5, enddate = $6, cabinet = $7
		WHERE idapplication = $8
	`
	_, err := db.Exec(query,
		application.Description.String,
		application.NameTeacher.String,
		application.IDEngineer.Int64,
		application.Status.String,
		application.StartDate,
		application.EndDate.Time,
		application.Cabinet.String,
		application.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи с ID %d: %v", application.ID, err)
	}
	return nil
}

func GenerateApplicationsPDF(db *sql.DB, filename string) error {
	query := `SELECT idapplication, description, nameteacher, idengineer, status, startdate, enddate, cabinet FROM applications`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("ошибка при запросе данных: %v", err)
	}
	defer rows.Close()

	var applications []types.Application
	for rows.Next() {
		var app types.Application
		if err := rows.Scan(&app.ID, &app.Description, &app.NameTeacher, &app.IDEngineer, &app.Status, &app.StartDate, &app.EndDate, &app.Cabinet); err != nil {
			return fmt.Errorf("ошибка при чтении данных: %v", err)
		}
		applications = append(applications, app)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("ошибка при обработке строк: %v", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddUTF8Font("Arial", "", "fonts/arialmt.ttf")
	pdf.SetFont("Arial", "", 12)

	pdf.AddPage()
	pdf.Cell(40, 10, "Список заявок")
	pdf.Ln(10)

	for _, app := range applications {
		formattedDate := app.StartDate.Format("2006-01-02")
		pdf.Cell(40, 10, fmt.Sprintf("ID: %d", app.ID))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Преподаватель: %s", app.NameTeacher.String))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Инженер ID: %d", app.IDEngineer.Int64))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Статус: %s", app.Status.String))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Начало: %s", formattedDate))
		pdf.Ln(5)
		if !app.EndDate.Time.IsZero() {
			pdf.Cell(40, 10, fmt.Sprintf("Конец: %s", app.EndDate.Time.Format("2006-01-02")))
		} else {
			pdf.Cell(40, 10, fmt.Sprintf("Конец: %s", "Нет"))
		}
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintf("Кабинет: %s", app.Cabinet.String))
		pdf.Ln(5)
		pdf.Cell(40, 10, fmt.Sprintln("================================================"))
		pdf.Ln(10)
	}

	if err := pdf.OutputFileAndClose(filename); err != nil {
		return fmt.Errorf("ошибка при сохранении PDF: %v", err)
	}

	return nil
}

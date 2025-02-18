package handlers

import (
	"att-diplom/internal/functions"
	"att-diplom/internal/helper"
	"att-diplom/internal/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	w.WriteHeader(http.StatusOK)
}

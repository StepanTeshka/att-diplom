package handlers

import (
	"att-diplom/internal/functions"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetApplicationsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	nameTeacher := r.URL.Query().Get("nameTeacher")
	engineerIDStr := r.URL.Query().Get("engineerID")
	orderDate := r.URL.Query().Get("orderDate")
	status := r.URL.Query().Get("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	var engineerID int
	if engineerIDStr != "" {
		var err error
		engineerID, err = strconv.Atoi(engineerIDStr)
		if err != nil {
			http.Error(w, "Неверный ID инженера", http.StatusBadRequest)
			return
		}
	}

	applications, err := functions.GetApplicationsWithPagination(db, page, pageSize, nameTeacher, engineerID, orderDate, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(applications)
	if err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"att-diplom/internal/functions"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetApplicationByIDHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("idApplication")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	application, err := functions.GetApplicationByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(application)
	if err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"att-diplom/internal/functions"
	"att-diplom/internal/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var application types.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err = functions.UpdateApplication(db, application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Задача с ID %d обновлена", application.ID)
}

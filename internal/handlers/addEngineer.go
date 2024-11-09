package handlers

import (
	"att-diplom/internal/functions"
	"att-diplom/internal/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddEngineerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var engineer types.Engineer
	err := json.NewDecoder(r.Body).Decode(&engineer)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	newID, err := functions.AddEngineer(db, engineer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Инженер добавлен с ID: %d", newID)
}

package handlers

import (
	"att-diplom/internal/functions"
	"att-diplom/internal/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdateEngineerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var engineer types.Engineer
	err := json.NewDecoder(r.Body).Decode(&engineer)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err = functions.UpdateEngineer(db, engineer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Инженер с ID %d обновлён", engineer.ID)
}

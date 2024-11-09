package handlers

import (
	"att-diplom/internal/functions"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func DeleteEngineerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = functions.DeleteEngineer(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Инженер с ID %d удалён", id)
}

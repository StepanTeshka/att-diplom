package handlers

import (
	"att-diplom/internal/functions"
	"database/sql"
	"net/http"
)

func GeneratePdfApplication(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := "applications.pdf"
		if err := functions.GenerateApplicationsPDF(db, filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=applications.pdf")
		http.ServeFile(w, r, filename)
	}
}

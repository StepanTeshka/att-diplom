package handlers

import (
	"att-diplom/internal/functions"
	"database/sql"
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var authReq AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// password := "123321"

	// // Генерация bcrypt-хеша с фактором стоимости 10
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Хешированный пароль:", string(hash))

	userID, err := functions.AuthenticateUser(db, authReq.Email, authReq.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := functions.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
	w.WriteHeader(http.StatusOK)
}

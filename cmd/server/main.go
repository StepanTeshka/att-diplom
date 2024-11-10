package main

import (
	"att-diplom/internal/app"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.RunBot()
	if err != nil {
		log.Fatalf("failed to run appBot: %s", err.Error())
	}

	err = app.RunSite()
	if err != nil {
		log.Fatalf("failed to run appSite: %s", err.Error())
	}

	// password := "12332199"
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Хеш пароля:", string(hash))

	log.Println("Запуск сервера на :5050...")
	err = http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}

func enableCors(w *http.ResponseWriter, origin string) {
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

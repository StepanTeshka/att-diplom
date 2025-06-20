package main

import (
	"att-diplom/internal/app"
	"att-diplom/internal/appinit"
	"context"
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	ctx := context.Background()

	appinit.InitDeps(ctx)

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	db = appinit.InitBD()

	go func() {
		err := a.RunBots(db)
		if err != nil {
			log.Fatalf("failed to run bots: %s", err.Error())
		}
	}()

	go func() {
		err := a.RunSite(db)
		if err != nil {
			log.Fatalf("failed to run site: %s", err.Error())
		}
	}()

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

// func enableCors(w *http.ResponseWriter, origin string) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", origin)
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
// 	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
// }

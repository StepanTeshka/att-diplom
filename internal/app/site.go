package app

import (
	"att-diplom/internal/appinit"
	"att-diplom/internal/handlers"
	"att-diplom/internal/middleware"
	"database/sql"
	"net/http"
)

var db *sql.DB

func InitApp() {
	db = appinit.InitBD()
}

func RunSite() error {
	InitApp()

	// Auth

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthHandler(w, r, db)
	})

	// Engineer

	http.HandleFunc("/getEngineers", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllEngineersHandler(w, r, db)
	}))

	http.HandleFunc("/getEngineer", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetEngineerByIDHandler(w, r, db)
	}))

	http.HandleFunc("/addEngineer", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.AddEngineerHandler(w, r, db)
	}))

	http.HandleFunc("/deleteEngineer", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteEngineerHandler(w, r, db)
	}))

	http.HandleFunc("/updateEngineer", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateEngineerHandler(w, r, db)
	}))

	// Application

	http.HandleFunc("/getApplications", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetApplicationsHandler(w, r, db)
	}))

	http.HandleFunc("/getApplication", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetApplicationByIDHandler(w, r, db)
	}))

	http.HandleFunc("/addApplication", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.AddApplicationHandler(w, r, db)
	}))

	http.HandleFunc("/deleteApplication", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteApplicationHandler(w, r, db)
	}))

	http.HandleFunc("/updateApplication", middleware.JwtAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateApplicationHandler(w, r, db)
	}))

	return nil
}

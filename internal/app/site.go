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

func RouteWrapper(path string, handlerFunc http.HandlerFunc, useJwtMiddleware bool) {
	finalHandler := handlerFunc

	if useJwtMiddleware {
		finalHandler = middleware.JwtAuthMiddleware(finalHandler)
	}

	finalHandler = middleware.CORSMiddleware(finalHandler)

	http.HandleFunc(path, finalHandler)
}

func RunSite() error {
	InitApp()

	// Auth routes
	RouteWrapper("/auth", func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthHandler(w, r, db)
	}, false)

	// Engineer routes
	RouteWrapper("/getEngineers", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllEngineersHandler(w, r, db)
	}, true)

	RouteWrapper("/getEngineer", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetEngineerByIDHandler(w, r, db)
	}, true)

	RouteWrapper("/addEngineer", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddEngineerHandler(w, r, db)
	}, true)

	RouteWrapper("/deleteEngineer", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteEngineerHandler(w, r, db)
	}, true)

	RouteWrapper("/updateEngineer", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateEngineerHandler(w, r, db)
	}, true)

	// Application routes
	RouteWrapper("/getApplications", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetApplicationsHandler(w, r, db)
	}, true)

	RouteWrapper("/getApplication", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetApplicationByIDHandler(w, r, db)
	}, true)

	RouteWrapper("/addApplication", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddApplicationHandler(w, r, db, BotTelegram)
	}, true)

	RouteWrapper("/deleteApplication", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteApplicationHandler(w, r, db)
	}, true)

	RouteWrapper("/updateApplication", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateApplicationHandler(w, r, db, BotTelegram)
	}, true)

	return nil
}

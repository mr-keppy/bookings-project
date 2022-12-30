package main

import (
	"github.com/mr-keppy/bookings/pkg/config"
	"github.com/mr-keppy/bookings/pkg/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


func routes(app *config.AppConfig) http.Handler{
	/*mux := pat.New()
	mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
	*/

	mux:=chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about",handlers.Repo.About)
	mux.Get("/generals-quarters",handlers.Repo.Generals)
	mux.Get("/majors-suite",handlers.Repo.Majors)
	mux.Get("/make-reservation",handlers.Repo.Reservations)
	mux.Get("/search-availability",handlers.Repo.Avilability)
	mux.Post("/search-availability",handlers.Repo.PostAvilability)
	mux.Get("/search-availability-json",handlers.Repo.PostAvilabilityJSON)
	mux.Get("/contacts",handlers.Repo.Contacts)

	fileServer:= http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static",fileServer))
	return mux
}
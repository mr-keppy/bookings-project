package main

import (
	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/handlers"
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
	mux.Post("/make-reservation",handlers.Repo.PostReservations)
	mux.Get("/search-availability",handlers.Repo.Avilability)
	mux.Post("/search-availability",handlers.Repo.PostAvilability)
	mux.Post("/search-availability-json",handlers.Repo.PostAvilabilityJSON)
	mux.Get("/contacts",handlers.Repo.Contacts)
	mux.Get("/reservation-summary",handlers.Repo.ReservationSummary)
	mux.Get("/choose-room/{id}",handlers.Repo.ChooseRoom)
	mux.Get("/book-room",handlers.Repo.BookRoom)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	fileServer:= http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static",fileServer))
	return mux
}
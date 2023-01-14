package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/driver"
	"github.com/mr-keppy/bookings/internal/models"
	"github.com/mr-keppy/bookings/internal/render"
)


var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{}
var pathToTemplates ="./../../templates"
var infoLog *log.Logger
var errorLog *log.Logger

func getRoutes() http.Handler{
	gob.Register(models.Reservation{})
	app.InProduction = false

	
	infoLog = log.New(os.Stdout, "INFO\t",log.Ldate | log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout,"ERROR\t",log.Ldate | log.Ltime| log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()

	session.Lifetime = 24*time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=kishorpadmanabhan password=")
	if err != nil {
		log.Fatal("Error while connecting db")
	}
	tc, err := CreateTestTemplateCache()

	if err != nil {
		log.Println("Error while creating template:CreateTestTemplateCache")
		log.Fatal(err)
	}

	log.Println(tc["home"])

	app.TemplateCache = tc
	app.UseCache = true
	repo := NewTestRepo(&app,db)
	NewHandler(repo)

	render.NewRenderer(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	mux:=chi.NewRouter()

	mux.Use(middleware.Recoverer)
	
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about",Repo.About)
	mux.Get("/generals-quarters",Repo.Generals)
	mux.Get("/majors-suite",Repo.Majors)
	mux.Get("/make-reservation",Repo.Reservations)
	mux.Post("/make-reservation",Repo.PostReservations)
	mux.Get("/search-availability",Repo.Avilability)
	mux.Post("/search-availability",Repo.PostAvilability)
	mux.Post("/search-availability-json",Repo.PostAvilabilityJSON)
	mux.Get("/contacts",Repo.Contacts)
	mux.Get("/reservation-summary",Repo.ReservationSummary)

	fileServer:= http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static",fileServer))
	return mux
}
func NoSurf(next http.Handler) http.Handler{
	csrfHandler:= nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}
func CreateTestTemplateCache()(map[string]*template.Template, error){
	myCache := map[string]*template.Template{}

	// get all of the files name from *page.temp from template
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if(err !=nil){
		return myCache, err
	}

	//range through all files
	for _, page:= range pages{
		name := filepath.Base(page)
		ts, err:= template.New(name).Funcs(functions).ParseFiles((page))
		if(err !=nil){
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if(err !=nil){
			return myCache, err
		}

		if(len(matches)>0){
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		}

		if(err !=nil){
			return myCache, err
		}

		myCache[name] = ts
	}
	return myCache, nil
}
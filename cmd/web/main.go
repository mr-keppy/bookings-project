package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/handlers"
	"github.com/mr-keppy/bookings/internal/helpers"
	"github.com/mr-keppy/bookings/internal/models"
	"github.com/mr-keppy/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//this is main
func run() error{
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
	
	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	return nil
}

func main() {

	err:= run()

	if(err != nil){
		log.Fatal(err)
	}
	// what going to store in session

	
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
	// fmt.Println((fmt.Sprintf("Starting applicaiton on port #:%d", portNumber)))
	//_ = http.ListenAndServe(portNumber, nil)
}

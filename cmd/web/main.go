package main

import (
	"github.com/mr-keppy/bookings/pkg/config"
	"github.com/mr-keppy/bookings/pkg/handlers"
	"github.com/mr-keppy/bookings/pkg/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

//this is main

func main() {

	app.InProduction = false

	session = scs.New()

	session.Lifetime = 24*time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	
	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = true
	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
	// fmt.Println((fmt.Sprintf("Starting applicaiton on port #:%d", portNumber)))
	//_ = http.ListenAndServe(portNumber, nil)
}

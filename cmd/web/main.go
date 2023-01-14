package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/driver"
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
func run() (*driver.DB, error){
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Room{})

	app.InProduction = false
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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

	//connect to db
	log.Println("connect to db")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=kishorpadmanabhan password=")
	if err != nil {
		log.Fatal("Error while connecting db")
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	return db, nil
}

func main() {

	db, err:= run()

	if(err != nil){
		log.Fatal(err)
	}
	// what going to store in session

	defer db.SQL.Close()
	defer  close(app.MailChan)

	listenForMail()

	//msg:= models.MailData{
	//	To: "kishor.padmanabhan@thoughtworks.com",
	//	From: "kishor338@gmail.com",
	//	Subject: "test email",
	//	Content: "Hello <b>Test</b>",
	//}
	//app.MailChan <- msg

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
	// fmt.Println((fmt.Sprintf("Starting applicaiton on port #:%d", portNumber)))
	//_ = http.ListenAndServe(portNumber, nil)

}

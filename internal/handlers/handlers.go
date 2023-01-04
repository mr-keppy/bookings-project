package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/driver"
	"github.com/mr-keppy/bookings/internal/forms"
	"github.com/mr-keppy/bookings/internal/helpers"
	"github.com/mr-keppy/bookings/internal/models"
	"github.com/mr-keppy/bookings/internal/render"
	"github.com/mr-keppy/bookings/internal/repository"
	"github.com/mr-keppy/bookings/internal/repository/dbrepo"
)

var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// Create new Repo
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {

	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL,a),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

// Home page is the home handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})

}

// Reservations page is the major handler
func (m *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// Post Reservations page is the major handler
func (m *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")


	layout := "2006-01-02"
	startDate, err:= time.Parse(layout,sd)
	if err!=nil{
		helpers.ServerError(w, err)
	}

	endDate, err:= time.Parse(layout,ed)
	if err!=nil{
		helpers.ServerError(w, err)
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err!=nil{
		helpers.ServerError(w, err)
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomId,
	}

	form := forms.New(r.PostForm)
	log.Println(reservation)
	form.Required("first_name","last_name","email","phone")
	form.MinLength("first_name",5,r)
	form.IsEmail("email")

	if !form.Valid(){
		data:= make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl",&models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newResId, err := m.DB.InsertReservation(reservation)



	if err !=nil{
		helpers.ServerError(w, err)
		return
	}

	restriction:= models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newResId,
	}

	_, err = m.DB.InsertRoomRestriction(restriction)

	if err !=nil{
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(),"reservation",reservation)
	
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}


// Contacts page is the major handler
func (m *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "contacts.page.tmpl", &models.TemplateData{})

}

// Avilability page is the major handler
func (m *Repository) Avilability(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})

}

// Post Avilability page is the major handler
func (m *Repository) PostAvilability(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("Search date is %s and end date is %s", start, end)))

}

type jsonResponse struct {
	OK      bool   `json:ok`
	Message string `json:"message"`
}

// Post Avilability page is the major handler
func (m *Repository) PostAvilabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "	")

	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Majors page is the major handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})

}

// Generals  page is the generals handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})

}

// About this about handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, World"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func( m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request){

	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)

	if !ok{
			//log.Println("cannot get item from session")
			// helpers.ServerError(w,errors.New("cannot get item from session"))
			m.App.Session.Put(r.Context(), "error","Can't get reservation from session")
			http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
			return
	}

	m.App.Session.Remove(r.Context(),"reservation")
	
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
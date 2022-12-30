package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mr-keppy/bookings/pkg/config"
	"github.com/mr-keppy/bookings/pkg/models"
	"github.com/mr-keppy/bookings/pkg/render"
)

var Repo *Repository

// Repository type
type Repository struct{

	App *config.AppConfig
}

// Create new Repo
func NewRepo(a *config.AppConfig) *Repository{

	return &Repository{
		App: a,
	}
}

func NewHandler(r *Repository){
	Repo = r
}
// Home page is the home handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})

}

// Reservations page is the major handler
func (m *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})

}

// Contacts page is the major handler
func (m *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w, r, "contacts.page.tmpl", &models.TemplateData{})

}

// Avilability page is the major handler
func (m *Repository) Avilability(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})

}
// Post Avilability page is the major handler
func (m *Repository) PostAvilability(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)
	start:= r.Form.Get("start_date")
	end:= r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("Search date is %s and end date is %s",start,end)))

}

type jsonResponse struct{
	OK bool `json:ok`
	Message string `json:"message"`
}

// Post Avilability page is the major handler
func (m *Repository) PostAvilabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK: true,
		Message: "Available",
	}
	
	out, err := json.MarshalIndent(resp,"","	")

	if err!=nil{
		log.Println(err)
	}

	log.Println(string(out))
	w.Header().Set("Content-Type","application/json")
	w.Write(out)

}

// Majors page is the major handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})

}
// Generals  page is the generals handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})

}

// About this about handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, World"

	remoteIP := m.App.Session.GetString(r.Context(),"remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

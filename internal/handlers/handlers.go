package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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
// Create new Test Repo
func NewTestRepo(a *config.AppConfig, db *driver.DB) *Repository {

	return &Repository{
		App: a,
		DB: dbrepo.NewTestingRepo(db.SQL,a),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

// login screen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
}
// login screen
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err!=nil{
		log.Println(err)
	}

	email:= r.Form.Get("email")
	password:= r.Form.Get("password")

	form := forms.New(r.PostForm)

	form.Required("email","password")

	if !form.Valid(){
		// to do
	}

	id, _, err := m.DB.Authenticate(email,password)
	if err!=nil{
		log.Println(err)
		m.App.Session.Put(r.Context(),"error","invalid user input")
		http.Redirect(w, r, "/user/login",http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(),"user_id",id)
	m.App.Session.Put(r.Context(),"flash","login successfully")
	http.Redirect(w, r, "/",http.StatusSeeOther)

}

// Home page is the home handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// Reservations page is the major handler
func (m *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	
	res, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	
	if !ok{
		helpers.ServerError(w,errors.New("cannot get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)

	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	log.Println("room name",room.RoomName)
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(),"reservation", res)
	sd:= res.StartDate.Format("2006-01-02")
	ed:= res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"]= ed

	log.Println(res)
	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})

}

// Post Reservations page is the major handler
func (m *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		helpers.ServerError(w,errors.New("cannot get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	

/*
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
*/
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

	//m.App.Session.Put(r.Context(),"reservation",reservation)

	restriction:= models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newResId,
		RestrictionID: 1,
	}

	_, err = m.DB.InsertRoomRestriction(restriction)

	if err !=nil{
		helpers.ServerError(w, err)
		return
	}

	//send notificaitons - first to guest

	htmlMessage := fmt.Sprintf(`
	<strong> Reservation Confirmation <strong><br>

	Dear %s:, <br>
	This is confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2000-01-01"), reservation.EndDate.Format("2000-01-01"))

	msg := models.MailData{
		To: reservation.Email,
		From: "me@here.com",
		Subject: "Reservation Confirmation",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg
	
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
// book room
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID,  _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")


	layout:= "2021-01-01"

	startDate,_ := time.Parse(layout,sd)
	endDate,_ := time.Parse(layout,ed)
	var res models.Reservation

	room, err := m.DB.GetRoomByID(res.RoomID)

	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(),"reservation", res)
	
	http.Redirect(w, r,"/make-reservation", http.StatusSeeOther)
	
}
// Post Avilability page is the major handler
func (m *Repository) PostAvilability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err:= time.Parse(layout,start)
	if err!=nil{
		helpers.ServerError(w, err)
		return
	}

	endDate, err:= time.Parse(layout,end)
	if err!=nil{
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate,endDate)
	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	if(len(rooms)==0){
		m.App.Session.Put(r.Context(),"error","No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data:= make(map[string]interface{})

	data["rooms"] = rooms

	res:= models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}
	
	m.App.Session.Put(r.Context(),"reservation",res)

	//for _, i:= range rooms{
	//	m.App.InfoLog.Println("Room:",i.ID, i.RoomName)
	//}

	//w.Write([]byte(fmt.Sprintf("Search date is %s and end date is %s", start, end)))

	render.Template(w, r, "choose-room.page.tmpl",&models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK      bool   `json:ok`
	Message string `json:"message"`
	RoomID string `json:"room_id`
	StartDate string `json:"start_date`
	EndDate string `json:"end_date`
}

// Post Avilability page is the major handler
func (m *Repository) PostAvilabilityJSON(w http.ResponseWriter, r *http.Request) {

	
	sd:= r.Form.Get("start_date")
	ed:= r.Form.Get("end_date")

	layout:= "2021-01-01"

	startDate,_ := time.Parse(layout,sd)
	endDate,_ := time.Parse(layout,ed)

	roomId,_ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomId)

	resp := jsonResponse{
		OK: available,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomId),
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
	sd:= reservation.StartDate.Format("2006-01-02")
	ed:= reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"]= ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request){
	roomId, err := strconv.Atoi(chi.URLParam(r,"id"))
	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	res, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomId

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
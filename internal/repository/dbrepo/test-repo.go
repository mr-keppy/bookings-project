package dbrepo

import (
	"time"

	"github.com/mr-keppy/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}
func (m *testDBRepo) InsertReservation(res models.Reservation) (int,error) {
	return 1, nil

}
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) (int, error){
	return 1, nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, room_id int) (bool, error){
	return true, nil
}

// return rooms available
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error){
	
	var rooms []models.Room
	return rooms, nil
}
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error){
	
	var room models.Room
	return room, nil
}
func (m *testDBRepo) GetUserByID(id int) (models.User, error){
	var u models.User
	return u, nil
}
//update user 
func (m *testDBRepo) UpdateUser(u models.User) ( error){
return nil
}

func (m *testDBRepo) Authenticate(email, password string) (int, string, error){

	return 1, "xxdfsd", nil
}
func (m *testDBRepo) AllReservations()([] models.Reservation, error) {
	var reservations [] models.Reservation
	return reservations, nil
}
func (m *testDBRepo) AllNewReservations()([] models.Reservation, error) {
	var reservations [] models.Reservation
	return reservations, nil
}
func (m *testDBRepo) GetReservationByID(id int)(models.Reservation, error) {
	var res models.Reservation
	return res, nil
}
func (m *testDBRepo) UpdateReservation(u models.Reservation) ( error){
	return nil
}
func (m *testDBRepo) DeleteReservation(id int) ( error){
	return nil
}
func (m *testDBRepo) UpdateProcessedForReservation(processed, id int) ( error){
	return nil
}
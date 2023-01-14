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
package repository

import (
	"time"

	"github.com/mr-keppy/bookings/internal/models"
)

type DatabaseRepo interface{
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction)(int, error)
	SearchAvailabilityByDatesByRoomID(start, end time.Time, room_id int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	Authenticate(email, password string) (int, string, error)
	UpdateUser(u models.User) ( error)
	AllReservations()([] models.Reservation, error)
	AllNewReservations()([] models.Reservation, error)
	GetReservationByID(id int)(models.Reservation, error)
	UpdateReservation(u models.Reservation) ( error)
	DeleteReservation(id int) ( error)
	UpdateProcessedForReservation(processed, id int) ( error)
}
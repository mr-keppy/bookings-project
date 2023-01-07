package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/mr-keppy/bookings/internal/models"
)

func (m *postgreDBRepo) AllUsers() bool {
	return true
}
func (m *postgreDBRepo) InsertReservation(res models.Reservation) (int,error) {


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	
	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil

}
func (m *postgreDBRepo) InsertRoomRestriction(res models.RoomRestriction) (int, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	
	var newID int

	stmt := `insert into room_restrictions (start_date, end_date, restriction_id, reservation_id, room_id, created_at, updated_at)
			values($1, $2, $3, $4, $5, $6, $7) returning id`

	err := m.DB.QueryRowContext(ctx, stmt, res.StartDate, res.EndDate,res.RestrictionID, res.ReservationID, res.RoomID, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (m *postgreDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, room_id int) (bool, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `
		select
			count(id)
		from
			room_restrictions
		where 
			$1 < end_date and $2 > start_date and room_id = $3;
	`

	row := m.DB.QueryRowContext(ctx,query, start, end, room_id)
	
	err:= row.Scan(&numRows)

	if err != nil {
		return false, err
	}
	if numRows==0 {
		return true, nil
	}
	return false, nil
}

// return rooms available
func (m *postgreDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select
			r.id, r.room_name
		from
			rooms r
		where r.id not in
		(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)
	`

	rows, err := m.DB.QueryContext(ctx,query, start, end)
	
	
	if err != nil {
		return rooms, err
	}

	for rows.Next(){
		var room models.Room
		err:= rows.Scan(
			&room.ID,
			&room.RoomName,

		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if  err = rows.Err(); err!=nil {
		return rooms, err
		
	}

	return rooms, nil
}
func (m *postgreDBRepo) GetRoomByID(id int) (models.Room, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at from rooms where id= $1
	`
	rows, err := m.DB.QueryContext(ctx,query, id)
	
	
	if err != nil {
		return room, err
	}

	for rows.Next(){
		
		err:= rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,

		)
		if err != nil {
			return room, err
		}
	}

	if  err = rows.Err(); err!=nil {
		return room, err
		
	}
	log.Println("room",room)
	return room, nil
}
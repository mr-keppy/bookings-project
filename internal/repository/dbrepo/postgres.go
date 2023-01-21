package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"
	"golang.org/x/crypto/bcrypt"

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
func (m *postgreDBRepo) GetUserByID(id int) (models.User, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	select id, first_name, last_name, email, password, access_level from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err:= row.Scan(
		&u.ID, u.FirstName, u.LastName, u.Email, u.Password, u.AccessLevel,
	)
	if err != nil{
		return u, err
	}else{
		return u, nil
	}
}

//update user 
func (m *postgreDBRepo) UpdateUser(u models.User) ( error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	update users set  first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5 where id = $6`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil{
		return err
	}else{
		return nil
	}
}

func (m *postgreDBRepo) Authenticate(email, password string) (int, string, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var id int
	var hashedPassword string

	row:= m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err:= row.Scan(&id, &hashedPassword)

	if err!=nil{
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword{
		return 0, "", errors.New("incorrect username/password")
	}else if err!=nil{
		return 0, "",err
	}

	return id, hashedPassword, nil
}
func (m *postgreDBRepo) AllReservations()([] models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations [] models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err !=nil{
		return reservations, err
	}

	defer rows.Close()

	for rows.Next(){
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err!=nil{
			return reservations, err
		}

		reservations = append(reservations, i)
	}
	if err = rows.Err(); err !=nil{
		return reservations, err 
	}
	return reservations, nil
}
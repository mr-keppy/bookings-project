package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5"

)

//db holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn  = &DB{}

const maOpenDBConn = 10

const maxIdleDBConn = 5

const maxDBLifeTime = 5 * time.Minute

// db holds the db connection pool
func ConnectSQL(dsn string) (*DB, error){
	d, err:= NewDatabase(dsn)
	if err !=nil{
		panic(err)
	}
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifeTime)
	d.SetMaxOpenConns(maOpenDBConn)

	dbConn.SQL = d

	err = testDB(d)

	if err!=nil{
		return  nil, err
	}

	return dbConn, nil
}

//try to ping db
func testDB(d *sql.DB) error{
	err:= d.Ping()
	if err!=nil{
		return err
	}
	return nil
}

// creates a new db for the app
func NewDatabase(dsn string)(*sql.DB, error){
	db, err := sql.Open("pgx", dsn)

	if err !=nil{
		return nil, err
	}

	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, nil
}
package dbrepo

import (
	"database/sql"

	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/repository"
)

type postgreDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}

type testDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}
func NewTestingRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
		DB: conn,
	}
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgreDBRepo{
		App: a,
		DB: conn,
	}
}

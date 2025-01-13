package handlers

import "database/sql"

type Handlers struct {
	Db *sql.DB
}

func New(db *sql.DB) *Handlers {
	return &Handlers{
		Db: db,
	}
}

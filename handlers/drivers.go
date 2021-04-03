package handlers

import (
	"database/sql"
	"log"
)

type DriverHandler struct {
	l  *log.Logger
	db *sql.DB
}

func NewDriverHandler(l *log.Logger, db *sql.DB) *DriverHandler {
	return &DriverHandler{l, db}
}

type KeyDriver struct{}

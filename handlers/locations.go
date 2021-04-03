package handlers

import (
	"database/sql"
	"log"
)

type LocationHandler struct {
	l  *log.Logger
	db *sql.DB
}

func NewLocationHandler(l *log.Logger, db *sql.DB) *LocationHandler {
	return &LocationHandler{l, db}
}

type KeyLocation struct{}

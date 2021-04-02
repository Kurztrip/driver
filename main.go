package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kurztrip/driver/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	addr     = ":8080"
	host     = "fullstack-postgres"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "driver_db"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	l := log.New(os.Stdout, "driver-api", log.LstdFlags)
	dh := handlers.NewDriverHandler(l, db)
	sm := mux.NewRouter()

	l.Println("Starting service in port", addr)

	//GET METHODS
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", dh.GetDrivers)
	getRouter.HandleFunc("/{id:[0-9]+}", dh.GetDriverWithID)

	//PUT METHODS
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", dh.UpdateDriver)
	putRouter.Use(dh.MiddlewareDriverValidation)

	//POST METHODS
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", dh.AddDriver)
	postRouter.Use(dh.MiddlewareDriverValidation)

	//DELETE METHODS
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", dh.DeleteDriver)

	s := http.Server{
		Addr:         addr,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	l.Println("closing database...")
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

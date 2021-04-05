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
	addr = ":8080"
	//host = "localhost"
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
	lh := handlers.NewLocationHandler(l, db)
	sm := mux.NewRouter()

	l.Println("Starting service in port", addr)

	//GET METHODS
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", handlers.GetService)
	getRouter.HandleFunc("/drivers", dh.GetDrivers)
	getRouter.HandleFunc("/drivers/{id:[0-9]+}", dh.GetDriverWithID)
	getRouter.HandleFunc("/locations", lh.GetLocations)
	getRouter.HandleFunc("/locations/{id:[0-9]+}", lh.GetLocationWithID)

	//PUT METHODS
	putDriverRouter := sm.Methods(http.MethodPut).Subrouter()
	putDriverRouter.HandleFunc("/drivers/{id:[0-9]+}", dh.UpdateDriver)
	putDriverRouter.Use(dh.MiddlewareDriverValidation)
	putLocationRouter := sm.Methods(http.MethodPut).Subrouter()
	putLocationRouter.HandleFunc("/locations/{id:[0-9]+}", lh.UpdateLocation)
	putLocationRouter.Use(lh.MiddlewareLocationValidation)

	//POST METHODS
	postDriverRouter := sm.Methods(http.MethodPost).Subrouter()
	postDriverRouter.HandleFunc("/drivers", dh.AddDriver)
	postDriverRouter.Use(dh.MiddlewareDriverValidation)
	postLocationRouter := sm.Methods(http.MethodPost).Subrouter()
	postLocationRouter.HandleFunc("/locations", lh.AddLocation)
	postLocationRouter.Use(lh.MiddlewareLocationValidation)

	//DELETE METHODS
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/drivers/{id:[0-9]+}", dh.DeleteDriver)
	deleteRouter.HandleFunc("/locations/{id:[0-9]+}", lh.DeleteLocation)

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

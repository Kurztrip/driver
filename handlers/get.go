package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kurztrip/driver/data"
	"github.com/gorilla/mux"
)

func GetService(rw http.ResponseWriter, h *http.Request) {
	e := json.NewEncoder(rw)
	e.Encode("Welcome to the driver microservice!")
}

func (dh *DriverHandler) GetDrivers(rw http.ResponseWriter, h *http.Request) {
	dh.l.Println("Handle GET drivers")
	e := json.NewEncoder(rw)
	var dvrs data.Drivers
	rows, err := dh.db.Query(`SELECT driver_id, driver_name, driver_surname, driver_age, driver_email, driver_address, driver_phone 
							FROM drivers LIMIT $1`, 100)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var dvr data.Driver
		err = rows.Scan(&dvr.ID, &dvr.Name, &dvr.Surname, &dvr.Age, &dvr.Email, &dvr.Address, &dvr.Phone)
		if err != nil {
			// handle this error
			panic(err)
		}
		dvrs = append(dvrs, &dvr)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	e.Encode(dvrs)
}

func (dh *DriverHandler) GetDriverWithID(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	e := json.NewEncoder(rw)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	dh.l.Println("Handle GET driver with  id:", id)

	sqlStatement := `SELECT driver_name, driver_surname, driver_age, driver_email, driver_address, driver_phone 
					FROM drivers 
					WHERE driver_id = $1;`

	dvr := data.Driver{ID: id}

	row := dh.db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&dvr.Name, &dvr.Surname, &dvr.Age, &dvr.Email, &dvr.Address, &dvr.Phone); err {
	case sql.ErrNoRows:
		e.Encode("No rows were returned!")
	case nil:
		e.Encode(dvr)
	default:
		panic(err)
	}
}

func (lh *LocationHandler) GetLocations(rw http.ResponseWriter, h *http.Request) {
	lh.l.Println("Handle GET locations")
	e := json.NewEncoder(rw)
	var locs data.Locations
	rows, err := lh.db.Query(`SELECT location_id, driver_id, latitude, longitude 
							FROM locations LIMIT $1`, 100)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var loc data.Location
		err = rows.Scan(&loc.ID, &loc.Driver_ID, &loc.Latitude, &loc.Longitude)
		if err != nil {
			// handle this error
			panic(err)
		}
		locs = append(locs, &loc)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	e.Encode(locs)
}

func (lh *LocationHandler) GetLocationWithID(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	e := json.NewEncoder(rw)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	lh.l.Println("Handle GET location with  id:", id)

	sqlStatement := `SELECT location_id, driver_id, latitude, longitude
					FROM locations 
					WHERE location_id = $1;`

	loc := data.Location{ID: id}

	row := lh.db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&loc.Driver_ID, &loc.Driver_ID, &loc.Latitude, &loc.Longitude); err {
	case sql.ErrNoRows:
		e.Encode("No rows were returned!")
	case nil:
		e.Encode(loc)
	default:
		panic(err)
	}
}

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

func (lh *LocationHandler) GetLocations(rw http.ResponseWriter, h *http.Request) {
	lh.l.Println("Handle GET locations")
	var locs data.Locations
	rows, err := lh.db.Query(`SELECT location_id, truck_id, latitude, longitude, location_time 
							FROM locations LIMIT $1`, 100)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var loc data.Location
		err = rows.Scan(&loc.ID, &loc.Truck_ID, &loc.Latitude, &loc.Longitude, &loc.Time)
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
	locs.ToJSON(rw)
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

	sqlStatement := `SELECT location_id, truck_id, latitude, longitude, location_time
					FROM locations 
					WHERE location_id = $1;`

	loc := data.Location{ID: id}

	row := lh.db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&loc.ID, &loc.Truck_ID, &loc.Latitude, &loc.Longitude, &loc.Time); err {
	case sql.ErrNoRows:
		e.Encode("No rows were returned!")
	case nil:
		loc.ToJSON(rw)
	default:
		panic(err)
	}
}

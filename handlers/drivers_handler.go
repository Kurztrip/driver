package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kurztrip/driver/data"
	"github.com/gorilla/mux"
)

type DriverHandler struct {
	l  *log.Logger
	db *sql.DB
}

func NewDriverHandler(l *log.Logger, db *sql.DB) *DriverHandler {
	return &DriverHandler{l, db}
}

func (d *DriverHandler) GetDrivers(rw http.ResponseWriter, h *http.Request) {
	d.l.Println("Handle GET drivers")
	e := json.NewEncoder(rw)
	var dvrs []data.Driver
	rows, err := d.db.Query(`SELECT driver_id, driver_name, driver_surname, driver_age, driver_email, driver_address, driver_phone 
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
		dvrs = append(dvrs, dvr)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	e.Encode(dvrs)
}

func (d *DriverHandler) GetDriverWithID(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	e := json.NewEncoder(rw)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	d.l.Println("Handle GET driver with  id:", id)

	sqlStatement := `SELECT driver_name, driver_surname, driver_age, driver_email, driver_address, driver_phone 
					FROM drivers 
					WHERE driver_id = $1;`

	dvr := data.Driver{ID: id}

	row := d.db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&dvr.Name, &dvr.Surname, &dvr.Age, &dvr.Email, &dvr.Address, &dvr.Phone); err {
	case sql.ErrNoRows:
		e.Encode("No rows were returned!")
	case nil:
		e.Encode(dvr)
	default:
		panic(err)
	}
}

func (d *DriverHandler) AddDriver(rw http.ResponseWriter, r *http.Request) {
	d.l.Println("Handle POST driver")

	sqlStatement := `INSERT INTO drivers (driver_name, driver_surname, driver_age, driver_email, driver_address, driver_phone)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING driver_id`

	dvr := r.Context().Value(KeyDriver{}).(*data.Driver)

	id := 0
	err := d.db.QueryRow(sqlStatement, dvr.Name, dvr.Surname, dvr.Age, dvr.Email, dvr.Address, dvr.Phone).Scan(&id)
	if err != nil {
		panic(err)
	}
	d.l.Println("Driver was added succesfully! id:", id)
}

func (d *DriverHandler) UpdateDriver(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	d.l.Println("Handle PUT driver ", id)

	sqlStatement := `UPDATE drivers
					SET driver_name = $2, driver_surname = $3, driver_age = $4, driver_email = $5, driver_address = $6, driver_phone = $7
					WHERE driver_id = $1;`

	dvr := r.Context().Value(KeyDriver{}).(*data.Driver)

	res, err := d.db.Exec(sqlStatement, id, dvr.Name, dvr.Surname, dvr.Age, dvr.Email, dvr.Address, dvr.Phone)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	d.l.Println("Rows Affected: ", count)
}

func (d *DriverHandler) DeleteDriver(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	d.l.Println("Handle DELETE driver ", id)

	sqlStatement := `DELETE FROM drivers WHERE driver_id = $1;`
	res, err := d.db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	d.l.Println("Rows Affected: ", count)
}

type KeyDriver struct{}

func (h *DriverHandler) MiddlewareDriverValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		dvr := &data.Driver{}

		err := dvr.FromJson(r.Body)
		if err != nil {
			h.l.Println("[ERROR] deserializing driver", err)
			http.Error(rw, "Error reading driver", http.StatusBadRequest)
			return
		}

		//validate driver
		err = dvr.Validate()
		if err != nil {
			h.l.Println("[ERROR] validating driver", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating driver: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		//add driver to context
		ctx := context.WithValue(r.Context(), KeyDriver{}, dvr)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

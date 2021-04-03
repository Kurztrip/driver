package handlers

import (
	"net/http"

	"github.com/Kurztrip/driver/data"
)

func (dh *DriverHandler) AddDriver(rw http.ResponseWriter, r *http.Request) {
	dh.l.Println("Handle POST driver")

	sqlStatement := `INSERT INTO drivers (driver_name, driver_surname, driver_age, driver_email, driver_address, driver_phone)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING driver_id`

	dvr := r.Context().Value(KeyDriver{}).(*data.Driver)

	id := 0
	err := dh.db.QueryRow(sqlStatement, dvr.Name, dvr.Surname, dvr.Age, dvr.Email, dvr.Address, dvr.Phone).Scan(&id)
	if err != nil {
		panic(err)
	}
	dh.l.Println("Driver was added succesfully! id:", id)
}

func (lh *LocationHandler) AddLocation(rw http.ResponseWriter, r *http.Request) {
	lh.l.Println("Handle POST location")

	sqlStatement := `INSERT INTO locations (driver_id, latitude, longitude)
					VALUES ($1, $2, $3)
					RETURNING location_id`

	loc := r.Context().Value(KeyLocation{}).(*data.Location)

	id := 0
	err := lh.db.QueryRow(sqlStatement, loc.Driver_ID, loc.Latitude, loc.Longitude).Scan(&id)
	if err != nil {
		panic(err)
	}
	lh.l.Println("Location was added succesfully! id:", id)
}

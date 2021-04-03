package handlers

import (
	"net/http"
	"strconv"

	"github.com/Kurztrip/driver/data"
	"github.com/gorilla/mux"
)

func (dh *DriverHandler) UpdateDriver(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	dh.l.Println("Handle PUT driver ", id)

	sqlStatement := `UPDATE drivers
					SET driver_name = $2, driver_surname = $3, driver_age = $4, driver_email = $5, driver_address = $6, driver_phone = $7
					WHERE driver_id = $1;`

	dvr := r.Context().Value(KeyDriver{}).(*data.Driver)

	res, err := dh.db.Exec(sqlStatement, id, dvr.Name, dvr.Surname, dvr.Age, dvr.Email, dvr.Address, dvr.Phone)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	dh.l.Println("Rows Affected: ", count)
}

func (lh *LocationHandler) UpdateLocation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	lh.l.Println("Handle PUT location ", id)

	sqlStatement := `UPDATE locations
					SET driver_id = $2, latitude = $3, longitude = $4
					WHERE location_id = $1;`

	loc := r.Context().Value(KeyLocation{}).(*data.Location)

	res, err := lh.db.Exec(sqlStatement, id, loc.Driver_ID, loc.Latitude, loc.Longitude)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	lh.l.Println("Rows Affected: ", count)
}

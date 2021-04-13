package handlers

import (
	"net/http"
	"strconv"

	"github.com/Kurztrip/driver/data"
	"github.com/gorilla/mux"
)

func (lh *LocationHandler) UpdateLocation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	lh.l.Println("Handle PUT location ", id)

	sqlStatement := `UPDATE locations
					SET truck_id = $2, latitude = $3, longitude = $4
					WHERE location_id = $1;`

	loc := r.Context().Value(KeyLocation{}).(*data.Location)

	res, err := lh.db.Exec(sqlStatement, id, loc.Truck_ID, loc.Latitude, loc.Longitude)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	lh.l.Println("Rows Affected: ", count)
}

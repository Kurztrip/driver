package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Kurztrip/driver/data"
)

func (lh *LocationHandler) AddLocation(rw http.ResponseWriter, r *http.Request) {
	lh.l.Println("Handle POST location")

	sqlStatement := `INSERT INTO locations (truck_id, latitude, longitude, location_time)
					VALUES ($1, $2, $3, $4)
					RETURNING location_id`

	loc := r.Context().Value(KeyLocation{}).(*data.Location)

	id := 0
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	err := lh.db.QueryRow(sqlStatement, loc.Truck_ID, loc.Latitude, loc.Longitude, formatted).Scan(&id)
	if err != nil {
		panic(err)
	}
	lh.l.Println("Location was added succesfully! id:", id)
}

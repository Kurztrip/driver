package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (lh *LocationHandler) DeleteLocation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", id)
		return
	}
	lh.l.Println("Handle DELETE location ", id)

	sqlStatement := `DELETE FROM locations WHERE location_id = $1;`
	res, err := lh.db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	lh.l.Println("Rows Affected: ", count)
}

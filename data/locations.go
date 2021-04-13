package data

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

type Location struct {
	ID        int       `json:"id"`
	Truck_ID  int       `json:"truck_id" validate:"required"`
	Latitude  float32   `json:"latitude" validate:"required,gte=-90,lte=90"`
	Longitude float32   `json:"longitude" validate:"required,gte=-180,lte=180"`
	Time      time.Time `json:"time"`
}

func (lc *Location) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(lc)
}

func (lc *Location) Validate() error {
	validate := validator.New()
	return validate.Struct(lc)
}

type Locations []*Location

func (lc *Location) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(lc)
}

func (lcs *Locations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(lcs)
}

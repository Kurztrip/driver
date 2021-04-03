package data

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type Driver struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Age     uint8  `json:"age" validate:"gte=0,lte=130"`
	Email   string `json:"email" validate:"required,email"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required,numeric"`
}

func (d *Driver) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

func (d *Driver) Validate() error {
	validate := validator.New()
	return validate.Struct(d)
}

type Drivers []*Driver

func (d *Drivers) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(d)
}

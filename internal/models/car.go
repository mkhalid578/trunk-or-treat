package models

import "time"

// create a struct for the car model
// the model will parse the json data from the request body and map it to the struct fields
type Car struct {
	ID                    int       `json:"id" db:"id"`
	Make                  string    `json:"make" db:"make"`
	Model                 string    `json:"model" db:"model"`
	Year                  int       `json:"model_year" db:"model_year"`
	Trim                  string    `json:"trim" db:"trim"`
	BodyStyle             string    `json:"body_style" db:"body_style"`
	Powertrain            string    `json:"powertrain" db:"powertrain"`
	CargoFullVolumeCuFt   *float64  `json:"cargo_vol_seats_folded_cu_ft" db:"cargo_vol_seats_folded_cu_ft"`
	CargoBehind2ndRowCuFt *float64  `json:"cargo_vol_behind_2nd_row_cu_ft" db:"cargo_vol_behind_2nd_row_cu_ft"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

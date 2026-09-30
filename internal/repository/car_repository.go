package repository

import (
	"car-api/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CarFilter struct {
	Year       *int
	Make       *string
	Model      *string
	Trim       *string
	BodyStyle  *string
	Powertrain *string
}

func AddCar(pool *pgxpool.Pool,
	model string,
	make string,
	year int,
	trim string,
	bodyStyle string,
	powertrain string) (*models.Car, error) {
	// Implementation for adding a car to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
		INSERT INTO cars (model, make, model_year, trim, body_style, powertrain) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, model, make, model_year, trim, body_style, powertrain`

	var car models.Car
	// Execute the query and scan the returned values into the car struct
	err := pool.QueryRow(ctx, query, model, make, year, trim, bodyStyle, powertrain).Scan(
		&car.ID,
		&car.Model,
		&car.Make,
		&car.Year,
		&car.Trim,
		&car.BodyStyle,
		&car.Powertrain)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

// get all cars
func GetAllCars(pool *pgxpool.Pool) ([]models.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
	SELECT id, model, make, model_year, trim, body_style, powertrain,
		cargo_vol_seats_folded_cu_ft,
		cargo_vol_behind_2nd_row_cu_ft,
		created_at
	FROM cars
	ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID,
			&car.Model,
			&car.Make,
			&car.Year,
			&car.Trim,
			&car.BodyStyle,
			&car.Powertrain,
			&car.CargoFullVolumeCuFt,
			&car.CargoBehind2ndRowCuFt,
			&car.CreatedAt)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}

// get car by id
func GetCarByID(pool *pgxpool.Pool, id int) (*models.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
	SELECT id, model, make, model_year, trim, body_style, powertrain, created_at 
	FROM cars
	WHERE id = $1
	`

	var car models.Car
	err := pool.QueryRow(ctx, query, id).Scan(
		&car.ID,
		&car.Model,
		&car.Make,
		&car.Year,
		&car.Trim,
		&car.BodyStyle,
		&car.Powertrain,
		&car.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

func UpdateCarVolume(pool *pgxpool.Pool,
	id int,
	cargoBehind2ndRowCuFt *float64,
	totalCargoVolumeCuFt *float64) (*models.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
		UPDATE cars
		SET cargo_vol_behind_2nd_row_cu_ft = COALESCE($1, cargo_vol_behind_2nd_row_cu_ft),
			cargo_vol_seats_folded_cu_ft = COALESCE($2, cargo_vol_seats_folded_cu_ft)
		WHERE id = $3
		RETURNING id, cargo_vol_behind_2nd_row_cu_ft, cargo_vol_seats_folded_cu_ft, created_at
		`
	var car models.Car
	err := pool.QueryRow(ctx, query, cargoBehind2ndRowCuFt, totalCargoVolumeCuFt, id).Scan(
		&car.ID,
		&car.CargoBehind2ndRowCuFt,
		&car.CargoFullVolumeCuFt,
		&car.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func DeleteCar(pool *pgxpool.Pool, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
		DELETE FROM cars
		WHERE id = $1
	`

	command, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return fmt.Errorf("car with id %d not found", id)
	}
	return nil
}

func GetCarsByYear(pool *pgxpool.Pool, year int) ([]models.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
		SELECT id, model, make, model_year, trim, body_style, powertrain, created_at 
		FROM cars
		WHERE model_year = $1
	`
	rows, err := pool.Query(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID,
			&car.Model,
			&car.Make,
			&car.Year,
			&car.Trim,
			&car.BodyStyle,
			&car.Powertrain,
			&car.CreatedAt)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

func GetCarsByFilter(pool *pgxpool.Pool, filter CarFilter) ([]models.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // release of the context resources

	query := `
		SELECT id, model, make, model_year, trim, body_style, powertrain,
			cargo_vol_seats_folded_cu_ft, cargo_vol_behind_2nd_row_cu_ft, created_at
		FROM cars
		WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.Year != nil {
		query += fmt.Sprintf(" AND model_year = $%d", argIndex)
		args = append(args, *filter.Year)
		argIndex++
	}
	if filter.Make != nil {
		query += fmt.Sprintf(" AND make = $%d", argIndex)
		args = append(args, *filter.Make)
		argIndex++
	}
	if filter.Model != nil {
		query += fmt.Sprintf(" AND model = $%d", argIndex)
		args = append(args, *filter.Model)
		argIndex++
	}
	if filter.Trim != nil {
		query += fmt.Sprintf(" AND trim = $%d", argIndex)
		args = append(args, *filter.Trim)
		argIndex++
	}
	if filter.BodyStyle != nil {
		query += fmt.Sprintf(" AND body_style = $%d", argIndex)
		args = append(args, *filter.BodyStyle)
		argIndex++
	}
	if filter.Powertrain != nil {
		query += fmt.Sprintf(" AND powertrain = $%d", argIndex)
		args = append(args, *filter.Powertrain)
		argIndex++
	}
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID,
			&car.Model,
			&car.Make,
			&car.Year,
			&car.Trim,
			&car.BodyStyle,
			&car.Powertrain,
			&car.CargoFullVolumeCuFt,
			&car.CargoBehind2ndRowCuFt,
			&car.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

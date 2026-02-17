package database

import (
	"database/sql"
	"mms/app/models"
)

func CreateCar(db *sql.DB, car *models.Car) error {
	query := `INSERT INTO cars (make, model, year, price, sold, client_id, number_plate, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			  RETURNING id, created_at, updated_at`

	return db.QueryRow(query, car.Make, car.Model, car.Year, car.Price, car.Sold, car.ClientID, car.NumberPlate).Scan(
		&car.ID, &car.CreatedAt, &car.UpdatedAt,
	)
}

func GetAllCars(db *sql.DB, search string, sold *bool, limit, offset int) ([]*models.Car, int, error) {
	var total int
	whereClause := "WHERE (make ILIKE $1 OR model ILIKE $1 OR number_plate ILIKE $1)"
	args := []interface{}{"%" + search + "%"}

	if sold != nil {
		whereClause += " AND sold = $2"
		args = append(args, *sold)
	}

	countQuery := `SELECT COUNT(*) FROM cars ` + whereClause
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	limitIdx := len(args) + 1
	offsetIdx := len(args) + 2
	query := `SELECT id, make, model, year, price, sold, client_id, number_plate, created_at, updated_at
			  FROM cars ` + whereClause + ` ORDER BY created_at DESC LIMIT $` + string(rune(48+limitIdx)) + ` OFFSET $` + string(rune(48+offsetIdx))

	args = append(args, limit, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cars []*models.Car
	for rows.Next() {
		c := &models.Car{}
		if err := rows.Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.Price, &c.Sold, &c.ClientID, &c.NumberPlate, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		cars = append(cars, c)
	}
	return cars, total, nil
}

func GetCarByID(db *sql.DB, id string) (*models.Car, error) {
	c := &models.Car{}
	query := `SELECT id, make, model, year, price, sold, client_id, number_plate, created_at, updated_at FROM cars WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.Price, &c.Sold, &c.ClientID, &c.NumberPlate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func UpdateCarSoldStatus(db *sql.DB, carID string, clientID string, sold bool, price float64) error {
	query := `UPDATE cars SET sold = $1, client_id = $2, price = $3, updated_at = NOW() WHERE id = $4`
	_, err := db.Exec(query, sold, clientID, price, carID)
	return err
}

func GetDashboardStats(db *sql.DB) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total Cars (not sold)
	var totalCars int
	err := db.QueryRow("SELECT COUNT(*) FROM cars WHERE sold = false").Scan(&totalCars)
	if err != nil {
		return nil, err
	}
	stats["total_cars"] = totalCars

	// Total Value (sum of prices of all cars)
	var totalValue sql.NullFloat64
	err = db.QueryRow("SELECT SUM(price) FROM cars").Scan(&totalValue)
	if err != nil {
		return nil, err
	}
	stats["total_value"] = totalValue.Float64

	// Total Profit (Since there is no complex transaction tracking yet,
	// we'll sum (price - expenses) for sold cars or just sum transactions)
	var totalProfit sql.NullFloat64
	err = db.QueryRow("SELECT SUM(amount) FROM transactions WHERE type = 'income'").Scan(&totalProfit)
	if err != nil {
		// If transactions table doesn't exist or fail, fallback to 0
		stats["total_profit"] = 0.0
	} else {
		stats["total_profit"] = totalProfit.Float64
	}

	// Total Debt (sum of remaining balance in debts)
	var totalDebt sql.NullFloat64
	err = db.QueryRow("SELECT SUM(remaining_balance) FROM debts WHERE status != 'paid'").Scan(&totalDebt)
	if err != nil {
		return nil, err
	}
	stats["total_debt"] = totalDebt.Float64

	return stats, nil
}

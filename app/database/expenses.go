package database

import (
	"database/sql"
	"mms/app/models"
)

func CreateExpense(db *sql.DB, expense *models.Expense) error {
	query := `INSERT INTO expenses (date, amount, category, description, car_id, client_id)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id, created_at`

	return db.QueryRow(query, expense.Date, expense.Amount, expense.Category, expense.Description, expense.CarID, expense.ClientID).Scan(
		&expense.ID, &expense.CreatedAt,
	)
}

func GetExpensesByCarID(db *sql.DB, carID string) ([]*models.Expense, error) {
	query := `SELECT e.id, e.date, e.amount, e.category, e.description, e.car_id, e.client_id, e.created_at,
					 COALESCE(c.number_plate, '') as car_plate
			  FROM expenses e
			  LEFT JOIN cars c ON e.car_id = c.id
			  WHERE e.car_id = $1 ORDER BY e.date DESC`

	rows, err := db.Query(query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		e := &models.Expense{}
		if err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.Category, &e.Description, &e.CarID, &e.ClientID, &e.CreatedAt, &e.CarPlate); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func GetAllExpenses(db *sql.DB, search string, limit, offset int) ([]*models.Expense, int, error) {
	searchQuery := "%" + search + "%"

	countQuery := `SELECT COUNT(*) FROM expenses 
				   WHERE category ILIKE $1 OR description ILIKE $1`
	var total int
	err := db.QueryRow(countQuery, searchQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT e.id, e.date, e.amount, e.category, e.description, e.car_id, e.client_id, e.created_at,
					 COALESCE(c.number_plate, '') as car_plate,
					 COALESCE(cl.name, '') as client_name
			  FROM expenses e
			  LEFT JOIN cars c ON e.car_id = c.id
			  LEFT JOIN clients cl ON e.client_id = cl.id
			  WHERE e.category ILIKE $1 OR e.description ILIKE $1 OR c.number_plate ILIKE $1 OR cl.name ILIKE $1
			  ORDER BY e.date DESC LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		e := &models.Expense{}
		if err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.Category, &e.Description, &e.CarID, &e.ClientID, &e.CreatedAt, &e.CarPlate, &e.ClientName); err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, e)
	}
	return expenses, total, nil
}

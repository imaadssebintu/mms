package database

import (
	"database/sql"
	"mms/app/models"
)

func CreateExpense(db *sql.DB, expense *models.Expense) error {
	query := `INSERT INTO expenses (date, amount, category, description, car_id)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, created_at`

	return db.QueryRow(query, expense.Date, expense.Amount, expense.Category, expense.Description, expense.CarID).Scan(
		&expense.ID, &expense.CreatedAt,
	)
}

func GetExpensesByCarID(db *sql.DB, carID string) ([]*models.Expense, error) {
	query := `SELECT id, date, amount, category, description, car_id, created_at
			  FROM expenses WHERE car_id = $1 ORDER BY date DESC`

	rows, err := db.Query(query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		e := &models.Expense{}
		if err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.Category, &e.Description, &e.CarID, &e.CreatedAt); err != nil {
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

	query := `SELECT id, date, amount, category, description, car_id, created_at
			  FROM expenses 
			  WHERE category ILIKE $1 OR description ILIKE $1
			  ORDER BY date DESC LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		e := &models.Expense{}
		if err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.Category, &e.Description, &e.CarID, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, e)
	}
	return expenses, total, nil
}

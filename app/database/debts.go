package database

import (
	"database/sql"
	"mms/app/models"
)

func CreateDebt(db *sql.DB, debt *models.Debt) error {
	query := `INSERT INTO debts (client_id, car_id, remaining_balance, payment_deadline, next_payment_date, status, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, NOW())
			  RETURNING id, created_at, updated_at`

	return db.QueryRow(query, debt.ClientID, debt.CarID, debt.RemainingBalance, debt.PaymentDeadline, debt.NextPaymentDate, debt.Status).Scan(
		&debt.ID, &debt.CreatedAt, &debt.UpdatedAt,
	)
}

func GetDebts(db *sql.DB, search string, limit, offset int) ([]interface{}, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM debts d
				   JOIN clients c ON d.client_id = c.id
				   JOIN cars cr ON d.car_id = cr.id
				   WHERE c.name ILIKE $1 OR cr.number_plate ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT d.id, d.client_id, d.car_id, d.remaining_balance, d.payment_deadline, d.next_payment_date, d.status, d.created_at, d.updated_at,
			  c.name as client_name, c.phone as client_phone, cr.number_plate as car_number_plate
			  FROM debts d
			  JOIN clients c ON d.client_id = c.id
			  JOIN cars cr ON d.car_id = cr.id
			  WHERE c.name ILIKE $1 OR cr.number_plate ILIKE $1
			  ORDER BY d.next_payment_date ASC LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var debts []interface{}
	for rows.Next() {
		var d models.Debt
		var clientName, clientPhone, carPlate string
		if err := rows.Scan(&d.ID, &d.ClientID, &d.CarID, &d.RemainingBalance, &d.PaymentDeadline, &d.NextPaymentDate, &d.Status, &d.CreatedAt, &d.UpdatedAt, &clientName, &clientPhone, &carPlate); err != nil {
			return nil, 0, err
		}

		debtMap := map[string]interface{}{
			"id":                d.ID,
			"client_id":         d.ClientID,
			"car_id":            d.CarID,
			"remaining_balance": d.RemainingBalance,
			"payment_deadline":  d.PaymentDeadline,
			"next_payment_date": d.NextPaymentDate,
			"status":            d.Status,
			"created_at":        d.CreatedAt,
			"updated_at":        d.UpdatedAt,
			"client_name":       clientName,
			"client_phone":      clientPhone,
			"car_number_plate":  carPlate,
		}
		debts = append(debts, debtMap)
	}
	return debts, total, nil
}

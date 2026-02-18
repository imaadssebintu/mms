package database

import (
	"database/sql"
	"mms/app/models"
)

func CreateClient(db *sql.DB, client *models.Client) error {
	query := `INSERT INTO clients (name, email, phone, address, notes, updated_at)
			  VALUES ($1, $2, $3, $4, $5, NOW())
			  RETURNING id, created_at, updated_at`

	return db.QueryRow(query, client.Name, client.Email, client.Phone, client.Address, client.Notes).Scan(
		&client.ID, &client.CreatedAt, &client.UpdatedAt,
	)
}

func GetAllClients(db *sql.DB, search string, limit, offset int) ([]*models.Client, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM clients WHERE name ILIKE $1 OR phone ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, email, phone, address, notes, created_at, updated_at
			  FROM clients WHERE name ILIKE $1 OR phone ILIKE $1
			  ORDER BY name ASC LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []*models.Client
	for rows.Next() {
		c := &models.Client{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		clients = append(clients, c)
	}
	return clients, total, nil
}

func GetClientByID(db *sql.DB, id string) (*models.Client, error) {
	c := &models.Client{}
	query := `SELECT id, name, email, phone, address, notes, created_at, updated_at FROM clients WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetClientByNameAndPhone(db *sql.DB, name, phone string) (*models.Client, error) {
	c := &models.Client{}
	query := `SELECT id, name, email, phone, address, notes, created_at, updated_at FROM clients WHERE name = $1 AND phone = $2 LIMIT 1`
	err := db.QueryRow(query, name, phone).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

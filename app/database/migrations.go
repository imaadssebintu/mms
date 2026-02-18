package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	queries := []string{
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"",

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password TEXT NOT NULL,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			email_confirmed BOOLEAN DEFAULT false,
			company_name VARCHAR(200),
			location VARCHAR(200),
			phone1 VARCHAR(20),
			phone2 VARCHAR(20),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS clients (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(200) NOT NULL,
			email VARCHAR(255),
			phone VARCHAR(20),
			address TEXT,
			notes TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS cars (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			make VARCHAR(100) NOT NULL,
			model VARCHAR(100) NOT NULL,
			year INTEGER,
			price DECIMAL(15, 2),
			sold BOOLEAN DEFAULT false,
			client_id UUID REFERENCES clients(id) ON DELETE SET NULL,
			number_plate VARCHAR(20) UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS expenses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			date TIMESTAMP WITH TIME ZONE NOT NULL,
			amount DECIMAL(15, 2) NOT NULL,
			category VARCHAR(100),
			description TEXT,
			car_id UUID REFERENCES cars(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS debts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
			car_id UUID NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
			remaining_balance DECIMAL(15, 2) NOT NULL,
			payment_deadline TIMESTAMP WITH TIME ZONE,
			next_payment_date TIMESTAMP WITH TIME ZONE,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type VARCHAR(50) NOT NULL,
			car_id UUID REFERENCES cars(id) ON DELETE SET NULL,
			client_id UUID REFERENCES clients(id) ON DELETE SET NULL,
			amount DECIMAL(15, 2) NOT NULL,
			date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS password_reset_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) NOT NULL,
			token VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			used BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Migration failed for query: %s\nError: %v", q, err)
			return err
		}
	}

	// Add new columns to cars table if they don't exist
	alterQueries := []string{
		"ALTER TABLE cars ADD COLUMN IF NOT EXISTS color VARCHAR(50)",
		"ALTER TABLE cars ADD COLUMN IF NOT EXISTS engine_number VARCHAR(100)",
		"ALTER TABLE cars ADD COLUMN IF NOT EXISTS chassis_number VARCHAR(100)",
		"ALTER TABLE cars ADD COLUMN IF NOT EXISTS purchase_price DECIMAL(15, 2)",
		"ALTER TABLE cars ADD COLUMN IF NOT EXISTS seller_id UUID",
		"ALTER TABLE cars DROP CONSTRAINT IF EXISTS cars_seller_id_fkey",
		"ALTER TABLE cars ADD CONSTRAINT cars_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES clients(id) ON DELETE SET NULL",
		"ALTER TABLE expenses ADD COLUMN IF NOT EXISTS client_id UUID REFERENCES clients(id) ON DELETE SET NULL",
	}

	for _, q := range alterQueries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Migration warning (alter): %v", err)
			// Don't fail matching logic, just log warning as column might exist
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

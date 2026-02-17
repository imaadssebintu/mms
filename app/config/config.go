package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	DB   *sql.DB
	SMTP SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var AppConfig *Config

func InitDB() {
	var psqlInfo string

	// Remote Database Connection Details
	host := "129.80.85.203"
	port := 5432
	user := "imaad"
	password := "Ertdfgxc"
	dbname := "mms"

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		log.Printf("Database connection failed: %v", err)
		log.Fatal("Cannot establish database connection")
	}

	AppConfig = &Config{
		DB: db,
		SMTP: SMTPConfig{
			Host:     "smtp.gmail.com",
			Port:     587,
			Username: "swadiqjuniorschools@gmail.com", // To be updated by user
			Password: "varn boqq brqq ftjv",           // To be updated by user
			From:     "swadiqjuniorschools@gmail.com",
		},
	}
	log.Println("Database connected successfully")
}

func GetDB() *sql.DB {
	return AppConfig.DB
}

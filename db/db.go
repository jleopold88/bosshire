package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "user=jamesleopold dbname=bosshire host=localhost sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(50) NOT NULL,
			role VARCHAR(10) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS jobs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			requirements TEXT NOT NULL,
			employer_id INTEGER REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS applications (
			id SERIAL PRIMARY KEY,
			job_id INTEGER REFERENCES jobs(id),
			talent_id INTEGER REFERENCES users(id),
			status VARCHAR(20) NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

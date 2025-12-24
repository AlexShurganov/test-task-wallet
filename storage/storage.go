package storage

import (
	"database/sql"
	"fmt"
	"log"
	"wallet-service/config"
)

func InitDB() (*sql.DB, error) {
	connStr := config.DBConnString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY,
			balance DECIMAL NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			wallet_id UUID NOT NULL,
			operation_type VARCHAR(10) NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
			amount DECIMAL NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (wallet_id) REFERENCES wallets(id)
		);`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	_, err := db.Exec(`INSERT INTO wallets (id, balance) VALUES ('f1b98811-f11f-44d4-a702-cd659b9e1c8c', 0) ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		log.Fatalf("Error creating test wallet")
	}

	_, err = db.Exec(`INSERT INTO wallets (id, balance) 
                 SELECT gen_random_uuid(), round((random() * 10000)::numeric, 2) 
                 FROM generate_series(1, 1000)`)
	if err != nil {
		log.Fatalf("Error creating wallets")
	}

	rows, err := db.Query("SELECT id FROM wallets ORDER BY RANDOM() LIMIT 10")
	if err != nil {
		log.Fatal("Error getting uuids")
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		rows.Scan(&id)
		fmt.Println(id)
	}
}

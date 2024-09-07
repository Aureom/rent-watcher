package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Init(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS properties (
            id TEXT PRIMARY KEY,
            first_photo TEXT,
            price TEXT,
            logradouro TEXT,
            bairro TEXT,
            cidade TEXT,
            metragem TEXT,
            quartos TEXT,
            banheiros TEXT,
            suites TEXT,
            garagens TEXT,
            tipo_imovel TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS raw_data (
            id TEXT PRIMARY KEY,
            json_data TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	return err
}

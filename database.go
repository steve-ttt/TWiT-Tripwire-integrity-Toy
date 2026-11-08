package tripwire

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func createDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS baseline (
			path TEXT PRIMARY KEY,
			hash TEXT NOT NULL
		);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveBaseline(dbPath string, data map[string]string) error {
	// check db file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Default().Println("Database file does not exist, creating a new one.")
		_, err := createDatabase(dbPath)
		if err != nil {
			return err
		}
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Default().Println("Failed to open database:", err)
		return err
	}
	for path, hash := range data {
		_, err := db.Exec("INSERT OR REPLACE INTO baseline (path, hash) VALUES (?, ?)", path, hash)
		if err != nil {
			log.Default().Println("Failed to insert/update baseline:", err)
			return err
		}
	}

	defer db.Close()
	return nil
}

func loadBaseline(dbPath string) (map[string]string, error) {
	data := make(map[string]string)
	// check db file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Default().Println("Database file does not exist.")
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Default().Println("Failed to open database:", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT path, hash FROM baseline")
	if err != nil {
		log.Default().Println("Failed to query baseline:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var path, hash string
		if err := rows.Scan(&path, &hash); err != nil {
			log.Default().Println("Failed to scan row:", err)
			return nil, err
		}
		data[path] = hash
	}

	if err := rows.Err(); err != nil {
		log.Default().Println("Row iteration error:", err)
		return nil, err
	}

	return data, nil
}

package db

const (
	CreateQuery = `
		CREATE TABLE IF NOT EXISTS router (
		id INTEGER NOT NULL PRIMARY KEY,
		source TEXT NOT NULL,
		target TEXT NOT NULL
		);
	`

	InsertSource = `
		
	`
)

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	ID        int64
	Username  string
	Type      string
	Content   string
	Timestamp time.Time
}

type DB struct {
	db *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_username ON events(username);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON events(timestamp);
	`

	_, err := db.Exec(query)
	return err
}

func (d *DB) SaveEvent(eventType, content string, timestamp time.Time, username string) error {
	query := `
	INSERT INTO events (type, content, timestamp, username)
	VALUES (?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, eventType, content, timestamp, username)
	return err
}

func (d *DB) GetEventsByUsername(username string) ([]Event, error) {
	query := `
	SELECT id, username, type, content, timestamp
	FROM events
	WHERE username = ?
	ORDER BY timestamp DESC
	`
	rows, err := d.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Username, &event.Type, &event.Content, &event.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (d *DB) GetAllUsernames() ([]string, error) {
	query := `
	SELECT DISTINCT username
	FROM events
	ORDER BY username
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, rows.Err()
}

func (d *DB) GetEventsByTimeRange(username string, start, end time.Time) ([]Event, error) {
	query := `
	SELECT id, username, type, content, timestamp
	FROM events
	WHERE username = ? AND timestamp BETWEEN ? AND ?
	ORDER BY timestamp DESC
	`
	rows, err := d.db.Query(query, username, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Username, &event.Type, &event.Content, &event.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (d *DB) DeleteOldEvents(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	query := `DELETE FROM events WHERE timestamp < ?`
	result, err := d.db.Exec(query, cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Export functions
func (d *DB) ExportToJSON(username string, outputPath string) error {
	events, err := d.GetEventsByUsername(username)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

func (d *DB) ExportToTXT(username string, outputPath string) error {
	events, err := d.GetEventsByUsername(username)
	if err != nil {
		return err
	}

	var content string
	for _, event := range events {
		content += fmt.Sprintf("[%s] %s: %s\n",
			event.Timestamp.Format("2006-01-02 15:04:05"),
			event.Type,
			event.Content)
	}

	return os.WriteFile(outputPath, []byte(content), 0644)
}

func (d *DB) ExportToDB(username string, outputPath string) error {
	// Create a new database for export
	exportDB, err := NewDB(outputPath)
	if err != nil {
		return err
	}
	defer exportDB.Close()

	// Get events for the username
	events, err := d.GetEventsByUsername(username)
	if err != nil {
		return err
	}

	// Insert events into the export database
	for _, event := range events {
		if err := exportDB.SaveEvent(event.Type, event.Content, event.Timestamp, event.Username); err != nil {
			return err
		}
	}

	return nil
} 
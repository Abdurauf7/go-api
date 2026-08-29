package models

import (
	"time"

	"api.com/api/db"
)

type Event struct {
	ID          int
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserId      int
}

// Save записывает событие в БД и проставляет сгенерированный id.
func (e *Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, datetime, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	return db.DB.QueryRow(query, e.Name, e.Description, e.Location, e.DateTime, e.UserId).Scan(&e.ID)
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT id, name, description, location, datetime, user_id FROM events`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}

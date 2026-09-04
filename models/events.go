package models

import (
	"database/sql"
	"errors"
	"time"

	"api.com/api/db"
)

// ErrEventNotFound возвращается, когда события с таким id нет в БД.
var ErrEventNotFound = errors.New("event not found")

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserId      int64
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
	query := `SELECT id, name, description, location, datetime, user_id FROM events ORDER BY id ASC`

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

func GetEventById(id int64) (*Event, error) {
	query := `SELECT * FROM events WHERE id = $1`

	row := db.DB.QueryRow(query, id)
	var e Event
	err := row.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventNotFound
	}

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (e *Event) Update() error {
	query := `
	UPDATE events 
	SET name = $2, description = $3, location = $4, datetime = $5
	WHERE id = $1`

	_, err := db.DB.Exec(query, e.ID, e.Name, e.Description, e.Location, e.DateTime)
	return err
}

func (e *Event) Delete() error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := db.DB.Exec(query, e.ID)
	return err
}

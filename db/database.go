package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Анонимный импорт драйвера (нужен только для инициализации)
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB — общий пул соединений, доступный всему приложению.
var DB *sql.DB

func InitDB() {
	// Строка подключения (URL-формат)
	// postgres://пользователь:пароль@хост:порт/название_бд
	dsn := "postgres://postgres:postgres@localhost:5432/events?sslmode=disable"

	// sql.Open не подключается к БД сразу, он лишь настраивает пул соединений
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Не удалось инициализировать БД: ", err)
	}

	// Настройки пула
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	// db.Ping физически устанавливает соединение и проверяет доступность базы
	if err := DB.Ping(); err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL: ", err)
	}

	fmt.Println("Успешно подключились к PostgreSQL!")

	createTables()
}

func createTables() {

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`

	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Fatal("Не удалось создать таблицу users: ", err)
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id          SERIAL PRIMARY KEY,
		name        TEXT NOT NULL,
		description TEXT NOT NULL,
		location    TEXT NOT NULL,
		datetime    TIMESTAMPTZ NOT NULL,
		user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
	)`

	if _, err := DB.Exec(createEventsTable); err != nil {
		log.Fatal("Не удалось создать таблицу events: ", err)
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations(
		id  SERIAL PRIMARY KEY, 
		event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE,
		user_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE (event_id, user_id)
	)`
	if _, err := DB.Exec(createRegistrationTable); err != nil {
		log.Fatal("Не удалось создать таблицу registrations: ", err)
	}
}

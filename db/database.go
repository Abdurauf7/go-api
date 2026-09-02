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
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id          SERIAL PRIMARY KEY,
		name        TEXT NOT NULL,
		description TEXT NOT NULL,
		location    TEXT NOT NULL,
		datetime    TIMESTAMPTZ NOT NULL,
		user_id     INTEGER NOT NULL
	)`

	if _, err := DB.Exec(createEventsTable); err != nil {
		log.Fatal("Не удалось создать таблицу events: ", err)
	}
}

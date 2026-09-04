package models

import (
	"errors"

	"api.com/api/db"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

// ErrUserExists возвращается, когда username уже занят.
var ErrUserExists = errors.New("user already exists")

// uniqueViolation — код ошибки Postgres при нарушении UNIQUE-ограничения.
const uniqueViolation = "23505"

type User struct {
	ID       int64
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// Save хеширует пароль, записывает пользователя в БД и проставляет сгенерированный id.
func (u *User) Save() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id`

	err = db.DB.QueryRow(query, u.Username, string(hashedPassword)).Scan(&u.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return ErrUserExists
	}

	return err
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id,username,password FROM users ORDER BY id ASC`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()

}
